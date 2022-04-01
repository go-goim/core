---
weight: 2
bookFlatSection: true
title: "实现异步并发 worker 队列"
date: 2022-04-01
tags:
    - "go"
    - "golang"
    - "worker"
    - "pool"
categories: 
    - "Development"
    - "golang"
---

> 记录实现一个异步并发 worker 队列的过程。

在开发 broadcast 功能的时候，碰到一个比较棘手的问题，需要并发执行多个 worker 来讲 broadcast 消息推送到所有在线用户，同时我希望能控制并发数量。

<!--more-->

## 前言

以往遇到类似的问题我都会借助 `sync.WaitGroup` 加 `channel` 的方式去做，实现方式也比较简单。大致思路如下：

```go
type LimitedWaitGroup struct {
    wg *sync.WaitGroup
    ch chan int
}

func NewLimitedWaitGroup(size int) *LimitedWaitGroup {
    return &LimitedWaitGroup{
        wg : new(sync.WaitGroup),
        ch : make(chan int, size)
    }
}

func (w *LimitedWaitGroup) Add(f func()) {
    // wait if channel is full
    w.ch <- 1
    w.wg.Add(1)
    go func() {
        defer w.done()
        f()
    }()
}

func (w *LimitedWaitGroup) done() {
    <-w.ch
    w.wg.Done()
}

func (w *LimitedWaitGroup) Wait() {
    w.wg.Wait()
}
```

这样能解决我大部分的简单需求，但是现在我想要的能力用这个简单的 `LimitedWaitGroup` 无法完全满足，所以重新设计了一个 `worker pool` 的概念来满足我现在以及以后类似的需求。

## 设计

### 需求整理

首先将目前想到的需求以及其优先级列出来：

**高优先级：**

1. worker pool 支持设置 size，防止 worker 无限增多
2. 任务并发执行且能指定并发数
3. 当 worker 达到上线时，新的任务在一定范围内支持排队等待（即 `limited queue`）
4. 支持捕获任务错误
5. 排队中的任务应该按顺序调度执行

**低优先级：**

1. 任务支持实时状态更新
2. 任务可以外部等待完成（类似 `waitGroup.Done()` ）
3. 当空闲 worker 小于指定并发数时，支持占用空闲 worker 部分运行（如当前剩余 3 个 worker 可用，但是新的任务需要 5 个并发，则尝试先占用这 3 个worker，并在运行过程中继续监听 pool 空闲出来的 worker 并尝试去占用）

{{< hint info >}}
**小结**  
列出完需求及其优先级后，经过考虑决定，高优先级除了`第五条`, 低优先级除了`第三条`, 其他需求都在目前版本里实现。

原因如下：

- 首先说低优先级第三条，这块的部分调度执行 worker，目前没有想好比较优雅的实现方式，所以暂时没有实现（但是下个版本会实现）
- 高优先级的第五条也是跟调度有点关系，如果队列里靠前的任务需要大量的 worker，那很容易造成阻塞，后面的 task 一直没办法执行，即便需要很少的 worker。所以等部分调度执行开发完再把任务按需执行打开。

{{< /hint >}}

### Task Definition

`task` 表示一次任务，包含了任务执行的方法，并发数，所属的 `workerSet`以及执行状态等。

```go
type TaskFunc func() error

type task struct {
    tf          TaskFunc        // task function
    concurrence int             // concurrence of task
    ws          *workerSet      // assign value after task distribute to worker.
    status      TaskStatus      // store task status.
}

// TaskStatus is the status of task.
type TaskStatus int
```

### Worker Definition

`worker` 作为最小调度单元，仅包含 `workerSet` 和 `error` .

```go
type worker struct {
    ws  *workerSet
    err error
}
```

### TaskResult Definition

`TaskResult` 是一个对外暴露的 `interface`, 用于外部调用者获取和管理任务执行状态信息。

```go
// TaskResult is a manager of submitted task.
type TaskResult interface {
    // get error if task failed.
    Err() error
    // wait for task done.
    Wait()
    // get task status.
    Status() TaskStatus
    // kill task.
    Kill()
}
```

{{< hint warning >}}
**`task` 和 `TaskStatus` 分别实现 `TaskResult` 的接口，从而外部统一拿到 `TaskResult`**

> 之所以 `TaskStatus` 也需要实现 `TaskResult` 是因为部分情况下，不需要创建 `task` 直接返回错误状态即可。如：
> 提交的任务的并发数过高（超过 pool 的 size），当前 queue 已满不能再处理任何其他任务了，这种情况直接返回对应的状态码。
{{< /hint >}}

### WorkerSet Definition

`workerSet` 为一组 `worker`的集合，作用是调度 `worker` 并维护起所属 `task` 的整个生命过程.

```go
// workerSet represent a group of task handle workers
type workerSet struct {
    task          *task
    runningWorker atomic.Int32
    workers       []*worker
    ctx           context.Context
    cancel        context.CancelFunc
    wg            *sync.WaitGroup
}
```

### WorkerPool Definition

`Pool` 是一个可指定 size 的 worker pool. 可并发运行多个 task 并且支持额外的任务排队能力。

```go
// Pool is a buffered worker pool
type Pool struct {
    // TODO: taskQueue should be a linked list, so that we can get the task from the head of the list and put it back to the head.
    // If we use a channel as taskQueue, we can't get the task from the head of the list and put it back to the head.
    // But make sure that before change it to linked list, we should have the ability run the task in min(taskQueue length, concurrence) goroutines.
    taskQueue         chan *task
    enqueuedTaskCount atomic.Int32 // count of unhandled tasks
    bufferSize        int          // size of taskQueue buffer, means can count of bufferSize task can wait to be handled
    maxWorker         int          // count of how many worker run in concurrence
    workerSets        []*workerSet
    lock              *sync.Mutex
    stopFlag          atomic.Bool
}
```

## 实现

上面已经确定需要的能力和基础的数据结构了，下面一个个去实现各个模块的能力。

### Worker Implement

`worker` 能力相对纯粹，看看 worker 是如何工作的:

```go
func (w *worker) run() {
    defer w.ws.done()

    var ec = make(chan error, 1)
    defer close(ec)
    go func() {
        ec <- w.ws.task.tf()
    }()

    select {
    case e := <-ec:
        w.err = e
    case <-w.ws.ctx.Done():
    }
}
```

### WorkerSet Implement

`workerSet` 调度 worker，记录 worker 运行状态等。

{{< details "点击展开" "...">}}

```go

func newWorkerSet(ctx context.Context, t *task) *workerSet {
    // 初始化参数
    // ... 省略代码
    return ws
}

func (ws *workerSet) run() {
    ws.task.updateStatus(TaskStatusRunning)
    for _, w := range ws.workers {
        ws.addOne()
        go w.run()
    }
}

func (ws *workerSet) stopAll() {
    ws.cancel()
    ws.task.updateStatus(TaskStatusKilled)
}

// err returns the first error that occurred in the workerSet.
func (ws *workerSet) err() error {
    // ...省略代码
    return nil
}

func (ws *workerSet) getRunningWorker() int {
    return int(ws.runningWorker.Load())
}

// done called when worker stop.
func (ws *workerSet) done() {
    ws.addRunningWorker(-1)
    ws.wg.Done()
    if ws.getRunningWorker() == 0 {
        ws.task.updateStatus(TaskStatusDone)
    }
}

// addOne called when worker start running.
func (ws *workerSet) addOne() {
    ws.addRunningWorker(1)
    ws.wg.Add(1)
}

func (ws *workerSet) wait() {
    ws.wg.Wait()
}
```

{{< /details >}}

### Task Implement

`task` 主要是记录 task 的状态，并通过 workerSet 控制其下的 worker.

{{< details "点击展开" "...">}}

```go
func newTask(tf TaskFunc, concurrence int) *task {
    return &task{
        tf:          tf,
        concurrence: concurrence,
    }
}

// Err returns the first error that occurred in the workerSet.
func (t *task) Err() error {
    // check t.ws if nil return nil.
    if t.ws == nil {
        return nil
    }

    return t.ws.err()
}

// Wait for task done.
// Please make sure task is done or running before call this function.
func (t *task) Wait() {
    // check t.ws if nil.
    if t.ws == nil {
        return
    }

    t.ws.wait()
}

// Status returns task status.
func (t *task) Status() TaskStatus {
    return t.status
}

// Kill task.
func (t *task) Kill() {
    // check t.ws if nil.
    if t.ws == nil {
        return
    }

    t.ws.stopAll()
}

func (t *task) assignWorkerSet(ws *workerSet) {
    t.ws = ws
}

func (t *task) updateStatus(status TaskStatus) {
    t.status = status
}
```

{{< /details >}}

### TaskStatus Implement

`TaskStatus` 虽然实现了 `TaskResult` 接口，但是不能控制任何 task，其有效的方法只有 `Status()` 和 `Err()`

```go
func (t TaskStatus) Error() string {
    return t.String()
}

func (t TaskStatus) Err() error {
    switch t {
    case TaskStatusError, TaskStatusQueueFull, TaskStatusTooManyWorker, TaskStatusPoolClosed, TaskStatusKilled:
        return t
    }

    return nil
}

func (t TaskStatus) Wait() {
    // do nothing.
}

func (t TaskStatus) Status() TaskStatus {
    return t
}

func (t TaskStatus) Kill() {
    // do nothing.
}
```

### Pool Implement

`Pool` 是总的入口，任务会提交到 `Pool`, 并由 `Pool` 创建 task 并调度到 `workerSet` 上，同时定时清理已完成的 `workerSet`, 
确保空闲 `worker` 能被合理使用。

{{< details "点击展开" "...">}}

```go
func NewPool(workerSize, queueSize int) *Pool {
    // ... 初始化各个参数

    // check p.enqueue to find out why make this channel size with p.bufferSize+1.
    //
    p.taskQueue = make(chan *task, p.bufferSize+1)
    // 启动单独 goroutine 维护队列
    go p.consumeQueue()
    return p
}

func (p *Pool) Submit(ctx context.Context, tf TaskFunc, concurrence int) TaskResult {
    if p.stopFlag.Load() {
        return TaskStatusPoolClosed
    }

    if concurrence > p.maxWorker {
        return TaskStatusTooManyWorker
    }

    // check if there has any worker place left
    p.lock.Lock()
    defer p.lock.Unlock()

    t := newTask(tf, concurrence)
    if p.tryRunTask(ctx, t) {
        return t
    }

    if p.enqueueTask(t, true) {
        return t
    }

    return TaskStatusQueueFull
}

func (p *Pool) Stop() {
    // 关闭队列和正在运行的 workerSet
}

// tryRunTask try to put task into workerSet and run it.Return false if capacity not enough.
// Make sure get p.Lock before call this func
func (p *Pool) tryRunTask(ctx context.Context, t *task) bool {
    if p.curRunningWorkerNum()+t.concurrence <= p.maxWorker {
        ws := newWorkerSet(ctx, t)
        p.workerSets = append(p.workerSets, ws)
        // run 为异步方法
        ws.run()
        return true
    }

    return false
}

// curRunningWorkerNum get current running worker num
// make sure lock mutex before call this func
func (p *Pool) curRunningWorkerNum() int {
    // ...省略代码
    return cnt
}

// enqueueTask put task to queue.
// p.enqueuedTaskCount increase 1 if is new task
func (p *Pool) enqueueTask(t *task, isNewTask bool) bool {
    // ... 省略代码
    return true
}

func (p *Pool) consumeQueue() {
    var ticker = time.NewTicker(time.Second)
    for {
        select {
        case t, ok := <-p.taskQueue:
            if !ok {
                // channel closed
                return
            }
            if p.tryRunTask(context.Background(), t) {
                p.enqueuedTaskCount.Sub(1)
                goto unlock
            }
            // if enqueueTask return false, means channel is closed.
            if !p.enqueueTask(t, false) {
                // channel is closed
                goto unlock
            }

        unlock:
            p.lock.Unlock()
        case <-ticker.C:
            log.Printf("current running worker num: %d", p.curRunningWorkerNum())
        }
    }

    // never reach here
}
```

{{< /details >}}

## 使用

到这里相关开发基本结束了，有一些 `TODO` 项后面后补充完善，下面通过 test case 来看一下如何使用这个 worker pool:


```go

func TestPool_SubmitOrEnqueue(t *testing.T) {
    p := NewPool(5, 1)
    var (
        cnt         int
        concurrence = 5
    )

    tf := func() error {
        time.Sleep(time.Second)
        log.Println("hello world")
        cnt++
        return nil
    }

    got := p.Submit(context.Background(), tf, concurrence)
    if got.Status() != TaskStatusRunning {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusRunning)
        return
    }
    got.Wait()
    if cnt != concurrence {
        t.Errorf("cnt = %v, want %v", cnt, concurrence)
    }
    if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusRunning {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got, TaskStatusRunning)
        return
    }
    if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusEnqueue {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusEnqueue)
        return
    }
    if got := p.Submit(context.Background(), tf, 6); got.Status() != TaskStatusTooManyWorker {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusTooManyWorker)
        return
    }
    if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusQueueFull {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusQueueFull)
        return
    }
    p.Stop()
    if got := p.Submit(context.Background(), tf, 1); got.Status() != TaskStatusPoolClosed {
        t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusPoolClosed)
        return
    }
}
```

## 总结

到这里这篇文章内容全部结束了，下面做一个简单的总结：

- 介绍背景和需求
- 根据需求定义了一组概念：`task`, `worker`, `workerSet`, `pool`
- 各个结构之前的关系以及如何实现
- 最终给出使用的 test case.

## 链接 🔗

**如果想仔细阅读源码，并持续关注这块功能的后续更新优化，请点击这里跳转到 [GitHub](https://github.com/yusank/goim/tree/main/pkg/worker).**
