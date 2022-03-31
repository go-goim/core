package pool

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Conn interface {
	Key() string           // connection unique key in pool
	Close() error          // close the connection
	Done() <-chan struct{} // check if ctx is canceled
	Err() error            // return error if conn has any internal error
}

type idleConn struct {
	p        *namedPool
	c        Conn
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
	var (
		timer = time.NewTimer(time.Second * 5)
	)
loop:
	for {

		select {
		case <-timer.C:
			timer.Reset(time.Second * 5)
			log.Infof("tick for conn=%s", i.c.Key())
			if i.c.Err() != nil {
				break loop
			}
		case <-i.stopChan:
			break loop
		case <-i.c.Done():
			log.Infof("conn done, err=%v", i.c.Err())
			break loop
		}
	}

	log.Infof("closing conn, key=%s", i.c.Key())
	i.close()
}
