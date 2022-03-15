package conn

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

type WsConn struct {
	Conn *websocket.Conn
	Key  string
}

type idleConn struct {
	*WsConn
	p        *pool
	stopChan chan struct{}
	t        time.Time
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
		case <-timer.C:
			break loop
		default:
			mt, message, err := i.Conn.ReadMessage()
			if err != nil {
				log.Info("read:", err)
				break loop
			}
			log.Infof("receiveType=%v, msg=%s", mt, message)
		}

	}

	log.Info("con-pool close conn")
	i.p.del(i.Key)
}

func (i *idleConn) stop() {
	select {
	case i.stopChan <- struct{}{}:
	default:
	}

	_ = i.Conn.Close()
}
