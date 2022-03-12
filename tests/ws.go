package tests

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func connectWs(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
