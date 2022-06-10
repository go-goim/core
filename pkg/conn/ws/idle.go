package ws

import (
	"github.com/go-goim/core/pkg/log"
)

type idleConn struct {
	*WebsocketConn
	p        *namedPool
	stopChan chan struct{}
}

// close is different form stop
// close is closes the connection and delete it from pool
// stop is a trigger to stop the daemon then call the close
func (i *idleConn) close() {
	_ = i.Close()
	i.p.delete(i.Key())
}

func (i *idleConn) stop() {
	i.stopChan <- struct{}{}
}

func (i *idleConn) daemon() {
loop:
	for {
		select {
		case <-i.ctx.Done():
			log.Error("conn done", "key", i.Key())
			break loop
		case <-i.stopChan:
			log.Info("conn stop", "key", i.Key())
			break loop
		case data := <-i.writeChan:
			i.writeToClient(data)
		}
	}

	log.Info("conn daemon exit", "key", i.Key())
	i.close()
}
