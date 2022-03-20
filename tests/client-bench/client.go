package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var clientMap = sync.Map{}

func loadClient(uid string) *client {
	v, ok := clientMap.Load(uid)
	if !ok {
		return nil
	}

	return v.(*client)
}

func initAndStoreClient(uid string) error {
	c, err := newClient(gatewayAddr, uid)
	if err != nil {
		return err
	}

	clientMap.Store(uid, c)
	return nil
}

type client struct {
	ws   *websocket.Conn
	uid  string
	stop chan struct{}
}

func newClient(addr, uid string) (*client, error) {
	pushAddr, err := getPushServerAddr(addr, uid)
	if err != nil {
		return nil, err
	}

	if pushAddr == "" {
		return nil, fmt.Errorf("push addr not found")
	}

	ws, err := connectWs(pushAddr, http.Header{"uid": []string{uid}})
	if err != nil {
		return nil, err
	}

	cli := &client{
		ws:   ws,
		uid:  uid,
		stop: make(chan struct{}, 1),
	}

	go cli.daemon()
	go cli.readMsgFromConn()
	return cli, nil
}

func (c *client) daemon() {
	var ticker = time.NewTimer(time.Second * 3)

loop:
	for {
		select {
		case <-ticker.C:
			err := c.ping()
			if err != nil {
				log.Println("ping err:", err)
				break loop
			}
		case <-c.stop:
			break loop
		}
	}
}

func (c *client) readMsgFromConn() {
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("read msg err:", err)
			return
		}
		log.Printf("Client=%s|data:%s\n", c.uid, string(data))
	}
}

func (c *client) ping() error {
	err := c.ws.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
	if err == nil {
		log.Printf("get user:%s ping", c.uid)
		return nil
	}

	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}

func getPushServerAddr(gatewayAddr, uid string) (addr string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/gateway/service/v1/discover", gatewayAddr), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("uid", uid)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code:%d,msg:%s", resp.StatusCode, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	defer resp.Body.Close()

	var result = make(map[string]string)
	if err = json.Unmarshal(b, &result); err != nil {
		return "", err
	}

	return result["agentId"], nil
}

func connectWs(addr string, h http.Header) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(addr, "http://"), Path: "/push/service/v1/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return nil, err
	}

	c.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
	return c, nil
}
