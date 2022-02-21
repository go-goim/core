package conn

import (
	"sync"

	"github.com/gorilla/websocket"
	messagev1 "github.com/yusank/goim/api/message/v1"
)

type WsConn struct {
	websocket.Conn
	mu       sync.Mutex
	pingChan <-chan struct{}
}

func (w *WsConn) PushMessage(message *messagev1.PushMessage) error {
	panic("implement me")
}

func (w *WsConn) Key() string {
	panic("implement me")
}

func (w *WsConn) Ping() <-chan struct{} {
	w.mu.Lock()
	if w.pingChan == nil {
		w.pingChan = make(chan struct{})
	}

	ch := w.pingChan
	w.mu.Unlock()
	return ch
}
