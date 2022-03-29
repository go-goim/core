package pool

import "sync"

var dp = newNamedPool()

func Add(c Conn) {
	dp.add(c)
}

func Get(key string) Conn {
	return dp.get(key)
}

func CloseAndDelete(key string) {
	dp.closeAndDelete(key)
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
	select {
	case <-c.Done():
		return
	default:
		if c.Err() != nil {
			return
		}
	}

	p.Lock()
	defer p.Unlock()
	i, loaded := p.m[c.Key()]
	if loaded {
		i.stop()
	}

	i = &idleConn{
		c:        c,
		stopChan: make(chan struct{}),
		p:        p,
	}

	go i.daemon()
	p.m[c.Key()] = i
}

func (p *namedPool) get(key string) Conn {
	p.RLock()
	i, ok := p.m[key]
	p.RUnlock()

	if ok {
		select {
		case <-i.c.Done():
			i.stop()
		default:
			if i.c.Err() != nil {
				i.stop()
				return nil
			}
			return i.c
		}
	}

	return nil
}

func (p *namedPool) closeAndDelete(key string) {
	p.RLock()
	i, ok := p.m[key]
	p.RUnlock()
	if !ok {
		return
	}

	// delete conn after close
	i.stop()
}

func (p *namedPool) delete(key string) {
	p.Lock()
	defer p.Unlock()
	delete(p.m, key)
}
