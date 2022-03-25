package pool

import "sync"

var dp = newNamedPool()

func Add(c Conn) {
	dp.add(c)
}

func Get(key string) Conn {
	return dp.get(key)
}

func Delete(key string) {
	dp.delete(key)
}

type namedPool struct {
	*sync.RWMutex
	m map[string]*idleConn
}

func newNamedPool() *namedPool {
	p := &namedPool{
		RWMutex: new(sync.RWMutex),
		m:       make(map[string]*idleConn),
	}

	return p
}

func (p *namedPool) add(c Conn) {
	if c.IsClosed() {
		return
	}

	p.Lock()
	defer p.Unlock()
	i, loaded := p.m[c.Key()]
	if loaded {
		i.stop()
	}

	i = &idleConn{
		c:        c,
		stopChan: make(chan struct{}, 1),
		p:        p,
	}

	go i.daemon()
	p.m[c.Key()] = i
}

func (p *namedPool) get(key string) Conn {
	p.RLock()
	i, ok := p.m[key]
	p.Unlock()

	if ok {
		if !i.c.IsClosed() {
			return i.c
		}

		go i.stop()
		p.delete(key)
	}

	return nil
}

func (p *namedPool) delete(key string) {
	p.Lock()
	defer p.Unlock()

	delete(p.m, key)
}
