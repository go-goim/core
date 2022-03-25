package pool

import "time"

type Conn interface {
	Key() string
	IsClosed() bool
	Close() error
	Reconcile() error
}

type idleConn struct {
	p        *namedPool
	c        Conn
	stopChan chan struct{}
}

// close is diffrent form stop
// close is close the connection and delete it from pool
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
		err   error
	)
loop:
	for {
		if err = i.c.Reconcile(); err != nil {
			break
		}

		select {
		case <-timer.C:
			timer.Reset(time.Second * 5)
		case <-i.stopChan:
			break loop
		}
	}

	i.close()
}
