package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/pkg/types"
	"github.com/go-goim/core/pkg/web/response"
)

type Friend struct {
	UID          types.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendUID    types.ID `json:"friendUid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendName   string   `json:"friendName" example:"friendName"`
	FriendAvatar string   `json:"friendAvatar" example:"https://www.example.com/friendAvatar.png"`
	// 0: friend, 1: stranger, 2: blacklist
	Status    int32 `json:"status" example:"0"`
	CreatedAt int64 `json:"createdAt" example:"1579098983"`
	UpdatedAt int64 `json:"updatedAt" example:"1579098983"`
}

type User struct {
	UID         types.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	Name        string   `json:"name" example:"user1"`
	Avatar      string   `json:"avatar" example:"https://www.example.com/avatar.png"`
	Email       *string  `json:"email,omitempty" example:"abc@example.com"`
	Phone       *string  `json:"phone,omitempty" example:"13800138000"`
	ConnectURL  *string  `json:"connectUrl,omitempty" example:"ws://10.0.0.1:8080/ws"`
	LoginStatus int32    `json:"loginStatus" example:"0"`
}

var (
	hostIPMode bool
	serverAddr string
	token      string
	logger     *log.Logger
	friends    []*Friend
	// user
	curUser  *User
	userName string
)

var (
	toUser = &Friend{}
)

const (
	loginURI = "/gateway/v1/user/login"
)

func init() {
	flag.BoolVar(&hostIPMode, "host_ip_mode", true, "use host ip instead of localhost")
	flag.StringVar(&serverAddr, "s", "127.0.0.1:18071", "gateway server addr")
	flag.StringVar(&userName, "u", "", "from user name")
	f, err := os.Create("./log.log")
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "[log]", log.Lshortfile)
	flag.Parse()
	if hostIPMode {
		serverAddr, err = getHostIP()
		assert(err == nil, "get host ip failed")
		serverAddr += ":18071"
	}
	userName += "@example.com"
}

func getHostIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// check the address type and do not use ipv6
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		// check the network type
		if ipnet.IP.IsLoopback() {
			continue
		}

		if ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("ip not found")
}

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

// TODO: support load friend list & support accept/reject friend request
//  remove logic of set toUid from flags,support select target user or parse target user form input message.

func main() {
	assert(userName != "", "from user id must be provided")

	addr, err := login()
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)

	conn, err := connectWs(addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err = loadFriends(); err != nil {
		panic(err)
	}

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
	if err = g.SetKeybinding("input", gocui.KeyCtrlK, gocui.ModNone, resetInput); err != nil {
		panic(err)
	}
	if err = g.SetKeybinding("input", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		panic(err)
	}

	// bind friends list view
	if err = g.SetKeybinding("friends", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		panic(err)
	}
	if err = g.SetKeybinding("friends", gocui.KeyArrowUp, gocui.ModNone, arrowUp); err != nil {
		panic(err)
	}
	if err = g.SetKeybinding("friends", gocui.KeyArrowDown, gocui.ModNone, arrowDown); err != nil {
		panic(err)
	}
	if err = g.SetKeybinding("friends", gocui.KeyEnter, gocui.ModNone, selectFriend); err != nil {
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

func connectWs(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(addr, "http://"), Path: "/push/v1/conn/ws"}
	h := http.Header{}
	h.Set("Authorization", token)
	logger.Println("token:", token)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return nil, err
	}

	c.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
	return c, nil
}

func login() (serverIP string, err error) {
	req := fmt.Sprintf(`{"email":"%s","password":"123456"}`, userName)

	resp, err := http.Post(fmt.Sprintf("http://%s%s", serverAddr, loginURI), "application/json", strings.NewReader(req))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		*response.BaseResponse
		Data *User
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if data.Code != 0 {
		return "", fmt.Errorf("login user=%s, err= %v", userName, data.Reason)
	}

	token = resp.Header.Get("Authorization")
	curUser = data.Data
	return *data.Data.ConnectURL, nil
}

func loadFriends() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/gateway/v1/user/friend/list", serverAddr), nil)
	if err != nil {
		logger.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println(rsp.StatusCode)
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Println(err)
		return err
	}

	var resp struct {
		response.BaseResponse
		Data []*Friend `json:"data"`
	}

	log.Println(string(b))
	err = json.Unmarshal(b, &resp)
	if err != nil {
		logger.Println(err)
		return err
	}

	if resp.Code != 0 {
		logger.Println(resp.Message)
		return fmt.Errorf(resp.Message)
	}

	friends = resp.Data
	return nil
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
			str := string(data)
			str = strings.Replace(str, "\n", "", -1)
			str = strings.Replace(str, "\r", "", -1)
			logger.Println("data:", str)
			dataChan <- data
		}
	}()

	return dataChan, errChan
}

func handleConn(conn *websocket.Conn, g *gocui.Gui, dataChan chan []byte) {
	for !layoutDone {
		time.Sleep(time.Millisecond * 100)
	}
	var (
		ticker = time.NewTicker(time.Second * 5)
	)
	for {
		select {
		case <-ticker.C:
			conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
		case data := <-dataChan:
			msg := new(messagev1.Message)
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
				fmt.Fprintf(v, "Receive|From:%d|Tp:%v|Content:%s|Seq:%d\n", msg.GetFrom(), msg.GetContentType(), msg.GetContent(), msg.GetMsgId())
				return nil
			})

		}
	}
}
