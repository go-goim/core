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

const (
	defaultSize    = 100
	defaultTimeout = time.Minute
)

var dp = newPool(defaultSize, defaultTimeout)

// GetConn return user connection
func GetConn(key string) (*WsConn, bool) {
	return dp.get(key)
}

// PutConn put the connection into pool
func PutConn(c *WsConn) error {
	return dp.put(c)
}

// RemoveConn close and remove a connection
func RemoveConn(key string) {
	dp.del(key)
}

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

func (p *pool) put(c *WsConn) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.connections) == p.size {
		return errPoolReachMaxSize
	}

	old, ok := p.connections[c.Key]
	if ok {
		old.stop()
	}

	ic := &idleConn{
		WsConn:   c,
		p:        p,
		stopChan: make(chan struct{}, 1),
		t:        time.Now(),
	}
	p.connections[c.Key] = ic

	go ic.daemon(p)
	return nil
}

func (p *pool) get(key string) (*WsConn, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	i, ok := p.connections[key]
	if ok {
		return i.WsConn, true
	}

	return nil, false
}

func (p *pool) del(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	i, ok := p.connections[key]
	if !ok {
		return
	}

	i.stop()
	delete(p.connections, key)
}
