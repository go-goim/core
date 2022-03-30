---
title: Introduction
type: docs
---

# GoIM

> Instant Messaging system written by Go.

## How to run

```shell
# run push server
make run Srv=push
# run gateway server
make run Srv=gateway
# run msg server
make run Srv=msg
```

## Design of GoIM

### 整体能力规划

![design](https://raw.githubusercontent.com/yusank/goim/main/static/images/goim.png)

#### 客户端如何查找和连接长连接服务

客户端如何连接长连接服务，目前我有两个方案各有优缺点，但是还没确定。

#### 反向代理方案

客户端统一入口在 gateway 上，gateway 支持反向代理能力，客户端发起长连接请求时，代理到后端的服务（这里准备使用一致性哈希来确定转发到哪台机器上）

优点:

- 入口统一，且可以在 gateway 上完成鉴权等操作
- 后端服务无需暴露 ip，且可任意扩缩容比较安全

缺点：

- gateway 需要承受长连接带来的压力，需要更多的 gateway 来承受大量在线用户的情况

#### httpdns 方案

客户端先通过暴露的域名，去访问 httpdns 服务获取真正后端服务的 ip，然后通过 ip 直接进行长连接

优点：

- 客户端与长连接服务器直连，减少代理层的压力

缺点：

- 要求暴露后端服务 ip，安全性降低且比较浪费 ip 资源

##### 纯 httpdns

![ws](https://raw.githubusercontent.com/yusank/goim/main/static/images/conn_ws_dns.png)

##### 结合 gateway

![gateway](https://raw.githubusercontent.com/yusank/goim/main/static/images/conn_ws_gateway.png)

#### 结论

最终决定,使用基于 gateway 作为第一入口,再返回长链接服务的方案.
原因如下:

1. 可以在 gateway 这一层做初步的校验和分配长链接服务的策略(比如按最小连接数,id 哈希等)
2. 反向代理会使系统更复杂且上层反向代理会有比较大的压力,项目初期不想搞太复杂
3. 对于客户端来说 gateway 就是一切了,之后要加的用户体系都是通过 gateway 暴露出来,入口可以比较收拢.

### 消息的流转

IM 数据将在 HBASE 上存储，关系型数据存在 MySQL

![msg](https://raw.githubusercontent.com/yusank/goim/main/static/images/send_rec_msg.png)
