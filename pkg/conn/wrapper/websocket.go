package wrapper

import (
	"context"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-goim/goim/pkg/log"
)

type WebsocketWrapper struct {
	context.Context
	*websocket.Conn
	UID    string
	cancel context.CancelFunc
	err    error
}

func WrapWs(ctx context.Context, c *websocket.Conn, uid string) *WebsocketWrapper {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx2, cancel := context.WithCancel(ctx)
	ww := &WebsocketWrapper{
		Context: ctx2,
		Conn:    c,
		UID:     uid,
		cancel:  cancel,
	}

	ww.SetCloseHandler(func(code int, text string) error {
		ww.cancelWithError(nil)
		message := websocket.FormatCloseMessage(code, "")
		_ = ww.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	ww.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == nil || err == websocket.ErrCloseSent {
			return nil
		}

		return err
	})

	return ww
}

func (w *WebsocketWrapper) AddCloseAction(f func() error) {
	cf := w.CloseHandler()
	w.SetCloseHandler(func(code int, text string) error {
		err := cf(code, text)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WebsocketWrapper) AddPingAction(f func() error) {
	pf := w.PingHandler()
	w.SetPingHandler(func(appData string) error {
		err := pf(appData)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WebsocketWrapper) cancelWithError(e error) {
	w.err = e
	w.cancel()
}

func (w *WebsocketWrapper) Key() string {
	return w.UID
}

func (w *WebsocketWrapper) Err() error {
	if w.err != nil {
		return w.err
	}

	if w.Context.Err() != nil {
		return w.Context.Err()
	}

	return nil
}

func (w *WebsocketWrapper) Close() error {
	// cancel context
	w.cancel()
	// close connection
	return w.Conn.Close()
}

// Daemon is keep read msg from connection, and handle registered ping, pong, close events
func (w *WebsocketWrapper) Daemon() {
	for {
		mt, message, err := w.ReadMessage()
		if err != nil {
			log.Error("websocket read message error", "error", err, "uid", w.UID)
			w.cancelWithError(err)
			return
		}
		log.Info("websocket read message", "uid", w.UID, "mt", mt, "message", string(message))
	}
}
