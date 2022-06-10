package ws

import "sync"

var dp = newNamedPool()

func addToPool(c *WebsocketConn) {
	dp.add(c)
}

func Get(key string) *WebsocketConn {
	return dp.get(key)
}

func LoadAllConn() chan *WebsocketConn {
	return dp.loadAllConns()
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

func (p *namedPool) add(c *WebsocketConn) {
	select {
	case <-c.ctx.Done():
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
		WebsocketConn: c,
		stopChan:      make(chan struct{}),
		p:             p,
	}

	go i.daemon()
	p.m[c.Key()] = i
}

func (p *namedPool) get(key string) *WebsocketConn {
	p.RLock()
	i, ok := p.m[key]
	p.RUnlock()

	if ok {
		select {
		case <-i.ctx.Done():
			i.stop()
		default:
			if i.Err() != nil {
				i.stop()
				return nil
			}
			return i.WebsocketConn
		}
	}

	return nil
}

func (p *namedPool) loadAllConns() chan *WebsocketConn {
	p.Lock()
	defer p.Unlock()

	ch := make(chan *WebsocketConn, len(p.m))
	for _, i := range p.m {
		ch <- i.WebsocketConn
	}

	close(ch)
	return ch
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
