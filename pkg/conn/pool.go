package conn

import (
	"errors"
	"sync"
	"time"
)

type pool struct {
	mu          sync.RWMutex
	connections map[string]*idleConn
	idleTimeout time.Duration
	size        int
}

type idleConn struct {
	conn     Conn
	stopChan chan struct{}
	t        time.Time
}

const (
	defaultSize    = 100
	defaultTimeout = time.Minute
)

// newPool 初始化连接
func newPool(size int, timeout time.Duration) *pool {
	if size <= 0 {
		size = defaultSize
	}

	if timeout <= 0 {
		timeout = defaultTimeout
	}

	p := &pool{
		connections: make(map[string]*idleConn, size),
		idleTimeout: timeout,
		size:        size,
	}

	return p
}

var errPoolReachMaxSize = errors.New("connection poll reach mas size")

func (p *pool) put(c Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.connections) == p.size {
		return errPoolReachMaxSize
	}

	old, ok := p.connections[c.Key()]
	if ok {
		old.stop()
	}

	ic := &idleConn{
		conn: c,
		t:    time.Now(),
	}
	p.connections[c.Key()] = ic

	if !ok {
		go ic.daemon(p)
	}

	return nil
}

func (i *idleConn) daemon(p *pool) {
	var (
		timer = time.NewTimer(p.idleTimeout)
	)
loop:
	for {
		select {
		case <-i.stopChan:
			break loop
		case <-i.conn.Ping():
			timer.Reset(p.idleTimeout)
		case <-timer.C:
			break loop
		}
	}

	_ = i.conn.Close()
}

func (i *idleConn) stop() {
	i.stopChan <- struct{}{}
}
