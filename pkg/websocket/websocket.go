package websocket

import (
	"net"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

type WrappedWs struct {
	*websocket.Conn
	closed bool
	UID    string
}

func WrapWs(c *websocket.Conn, uid string) *WrappedWs {
	ww := &WrappedWs{
		Conn:   c,
		UID:    uid,
		closed: false,
	}

	ww.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		_ = ww.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		ww.closed = true
		return nil
	})

	ww.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == nil || err == websocket.ErrCloseSent {
			return nil
		}

		if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	return ww
}

func (w *WrappedWs) AddCloseAction(f func() error) {
	cf := w.CloseHandler()
	w.SetCloseHandler(func(code int, text string) error {
		err := cf(code, text)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WrappedWs) AddPingAction(f func() error) {
	pf := w.PingHandler()
	w.SetPingHandler(func(appData string) error {
		err := pf(appData)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WrappedWs) Key() string {
	return w.UID
}

func (w *WrappedWs) IsClosed() bool {
	return w.closed
}

func (w *WrappedWs) Reconcile() error {
	mt, message, err := w.ReadMessage()
	if err != nil {
		log.Info("read:", err)
		return err
	}
	log.Infof("receiveType=%v, msg=%s", mt, message)
	return nil
}
