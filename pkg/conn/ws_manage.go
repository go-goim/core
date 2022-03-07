package conn

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	messagev1 "github.com/yusank/goim/api/message/v1"
)

type WsConn struct {
	*websocket.Conn
	key      string
	pingChan chan struct{}
}

func (wc *WsConn) PushMessage(message *messagev1.PushMessageReq) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_ = wc.SetWriteDeadline(time.Now().Add(time.Second))
	w, err := wc.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return w.Close()
}

func (w *WsConn) Key() string {
	return w.key
}

func (w *WsConn) Ping() <-chan struct{} {
	return w.pingChan
}

func (w *WsConn) pingFunc(message string) error {
	err := w.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	if err == nil {
		select {
		// try put ping
		case w.pingChan <- struct{}{}:
		// non-blocking
		default:
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

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1 << 16,
	ReadBufferSize:  1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	//todo use check uid/token middleware before this handler
	uid := r.Header.Get("uid")
	if uid == "" {
		log.Println("uid not found")
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// todo
		log.Println(err)
		return
	}

	wc := &WsConn{
		Conn:     c,
		pingChan: make(chan struct{}, 1),
	}

	wc.SetPingHandler(wc.pingFunc)
	_ = dp.put(wc)
}
