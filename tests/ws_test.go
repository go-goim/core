package tests

import (
	"testing"

	"golang.org/x/net/websocket"
)

func Test_connectWs(t *testing.T) {
	var addr string
	c, err := connectWs(addr)
	if err != nil {
		t.Error(err)
		return
	}

	if err = c.WriteMessage(websocket.TextFrame, []byte("hello")); err != nil {
		t.Error(err)
	}
}
