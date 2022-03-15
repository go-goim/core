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

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
)

var (
	serverAddr string
	uid        string
	logger     *log.Logger
)

func init() {
	flag.StringVar(&serverAddr, "s", "127.0.0.1:8071", "gateway server addr")
	flag.StringVar(&uid, "u", "user1", "user id")
	f, err := os.Create("log.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "[log]", log.Lshortfile)
	flag.Parse()
}

func main() {
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

	if err = g.MainLoop(); err != nil {
		fmt.Println("exit:", err)
	}
}

func layout(g *gocui.Gui) error {
	var views = []string{outputView, inputView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		logger.Println(x0, y0, x1, y1)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			logger.Println(err)
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + view + " "

			if view == inputView {
				v.Editable = true
				v.Wrap = true
			}

			if err != gocui.ErrUnknownView {
				return err
			}
		}
	}

	_, err := g.SetCurrentView(inputView)
	if err != nil {
		log.Fatal("failed to set current view: ", err)
	}
	return nil

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

	return result["agentID"], nil
}
