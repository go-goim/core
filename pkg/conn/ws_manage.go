package conn

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	messagev1 "github.com/yusank/goim/api/message/v1"
)

type WsConn struct {
	*websocket.Conn
	key      string
	mu       sync.Mutex
	pingChan chan struct{}
}

func (w *WsConn) PushMessage(message *messagev1.PushMessage) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return w.WriteMessage(websocket.TextMessage, b)
}

func (w *WsConn) Key() string {
	return w.key
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

// closedChan is a reusable closed channel.
var closedChan = make(chan struct{})

func (w *WsConn) pingFunc(message string) error {
	err := w.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	if err == nil {
		if w.pingChan == nil {
			w.pingChan = closedChan
		} else {
			close(w.pingChan)
		}

		return nil
	}

	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}
