package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"

	messagev1 "github.com/yusank/goim/api/message/v1"
)

var (
	serverAddr string
	uid        string
	toUid      string
	logger     *log.Logger
)

func init() {
	flag.StringVar(&serverAddr, "s", "127.0.0.1:8071", "gateway server addr")
	flag.StringVar(&uid, "u", "", "from user id")
	flag.StringVar(&toUid, "t", "", "to user id")
	f, err := os.Create("log.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "[log]", log.Lshortfile)
	u, _ := url.Parse("discovery://goim.push.service")
	logger.Println(u.Scheme)
	logger.Println(u.Host)
	logger.Println(u.Path)
	logger.Println(u.Opaque)
	flag.Parse()
}
func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func main() {
	assert(uid != "", "from user id must be provided")
	assert(toUid != "", "to user id must be provided")
	assert(uid != toUid, "uid and toUid must be different")

	addr, err := getPushServerAddr()
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
	conn, err := connectWs(addr, http.Header{"uid": []string{uid}})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	g.Cursor = true

	g.SetManagerFunc(layout)
	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err = g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, resetInput); err != nil {
		panic(err)
	}

	dc, ec := readMsgFromConn(conn)
	go handleConn(conn, g, dc)

	go func() {
		ec <- g.MainLoop()
	}()

	if err = <-ec; err != nil {
		g.Close()
		fmt.Println("exit:", err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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

func getPushServerAddr() (addr string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/gateway/service/v1/discover", serverAddr), nil)
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

func readMsgFromConn(conn *websocket.Conn) (chan []byte, chan error) {
	var (
		dataChan = make(chan []byte, 1)
		errChan  = make(chan error, 1)
	)

	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				logger.Println(err)
				errChan <- err
				return
			}
			logger.Println("data:", string(data))
			dataChan <- data
		}
	}()

	return dataChan, errChan
}

func handleConn(conn *websocket.Conn, g *gocui.Gui, dataChan chan []byte) {
	var (
		ticker = time.NewTicker(time.Second * 5)
	)
	for {
		select {
		case <-ticker.C:
			conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
		case data := <-dataChan:
			msg := new(messagev1.SendMessageReq)
			if err := json.Unmarshal(data, msg); err != nil {
				logger.Println("unmarshal err:", err)
				msg.Content = string(data)
				msg.ContentType = messagev1.MessageContentType_Text
			}

			g.Update(func(gg *gocui.Gui) error {
				v, err1 := gg.View("output")
				if err1 != nil {
					logger.Println("update err:", err1)
					return err1
				}
				fmt.Fprintln(v, "------")
				fmt.Fprintf(v, "Receive|From:%s|Tp:%v|Content:%s\n", msg.GetFromUser(), msg.GetContentType(), msg.GetContent())
				return nil
			})

		}
	}
}
