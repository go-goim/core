package ws

import (
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-goim/core/pkg/log"
)

type WebsocketConn struct {
	*websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc

	uid          string
	writeChan    chan []byte
	onWriteError func()
	err          error
}

var (
	ErrWriteChanFull = errors.New("write chan full")
)

func WrapWs(ctx context.Context, c *websocket.Conn, uid string) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx2, cancel := context.WithCancel(ctx)
	wc := &WebsocketConn{
		ctx:       ctx2,
		Conn:      c,
		uid:       uid,
		writeChan: make(chan []byte, 1),
		cancel:    cancel,
	}

	wc.SetCloseHandler(func(code int, text string) error {
		wc.cancelWithError(nil)
		message := websocket.FormatCloseMessage(code, "")
		_ = wc.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	wc.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		if err == nil || err == websocket.ErrCloseSent {
			return nil
		}

		return err
	})

	go wc.readDaemon()
	// add to pool
	addToPool(wc)
}

func (w *WebsocketConn) AddCloseAction(f func() error) {
	cf := w.CloseHandler()
	w.SetCloseHandler(func(code int, text string) error {
		err := cf(code, text)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WebsocketConn) AddPingAction(f func() error) {
	pf := w.PingHandler()
	w.SetPingHandler(func(appData string) error {
		err := pf(appData)
		if err == nil {
			return f()
		}

		return err
	})
}

func (w *WebsocketConn) cancelWithError(e error) {
	w.err = e
	w.cancel()
}

func (w *WebsocketConn) Key() string {
	return w.uid
}

func (w *WebsocketConn) Err() error {
	if w.err != nil {
		return w.err
	}

	if w.ctx.Err() != nil {
		return w.ctx.Err()
	}

	return nil
}

func (w *WebsocketConn) Close() error {
	// cancel context
	w.cancel()
	// close connection
	return w.Conn.Close()
}

// readDaemon is keep read msg from connection, and handle registered ping, pong, close events
func (w *WebsocketConn) readDaemon() {
	for {
		mt, message, err := w.ReadMessage()
		if err != nil {
			log.Error("websocket read message error", "error", err, "uid", w.uid)
			w.cancelWithError(err)
			return
		}
		log.Info("websocket read message", "uid", w.uid, "mt", mt, "message", string(message))
	}
}

func (w *WebsocketConn) Write(data []byte) error {
	select {
	case w.writeChan <- data:
		return nil
	default:
	}

	timer := time.NewTimer(time.Millisecond * 500)
	select {
	case w.writeChan <- data:
		return nil
	case <-timer.C:
		return ErrWriteChanFull
	}
}

func (w *WebsocketConn) write(data []byte) {
	_ = w.SetWriteDeadline(time.Now().Add(time.Second))
	err := w.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		w.onWriteError()
		return
	}
}
