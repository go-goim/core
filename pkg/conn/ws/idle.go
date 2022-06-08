package ws

import (
	"github.com/go-goim/core/pkg/log"
)

type idleConn struct {
	p        *namedPool
	c        *WebsocketConn
	stopChan chan struct{}
}

// close is different form stop
// close is closes the connection and delete it from pool
// stop is a trigger to stop the daemon then call the close
func (i *idleConn) close() {
	_ = i.c.Close()
	i.p.delete(i.c.Key())
}

func (i *idleConn) stop() {
	i.stopChan <- struct{}{}
}

func (i *idleConn) daemon() {
loop:
	for {
		select {
		case <-i.c.ctx.Done():
			log.Error("conn done", "key", i.c.Key())
			break loop
		case <-i.stopChan:
			log.Info("conn stop", "key", i.c.Key())
			break loop
		case data := <-i.c.writeChan:
			i.c.write(data)
		}
	}

	log.Info("conn daemon exit", "key", i.c.Key())
	i.close()
}
