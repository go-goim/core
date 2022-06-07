package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jroimartin/gocui"

	messagev1 "github.com/go-goim/api/message/v1"
)

func layout(g *gocui.Gui) error {
	var views = []string{outputView, inputView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		//logger.Println(x0, y0, x1, y1)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			logger.Println(err)
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + toUid + " "

			if view == inputView {
				v.Editable = true
				v.Wrap = true
				v.Title = " " + uid + " "
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

func resetInput(g *gocui.Gui, v *gocui.View) error {
	buf := &bytes.Buffer{}
	io.Copy(buf, v)
	// todo need load friend list then send msg
	m := &messagev1.SendMessageReq{
		FromUser:    "",
		ToUser:      "",
		ContentType: 0,
		Content:     "",
	}
	b, err := json.Marshal(&m)
	if err != nil {
		logger.Println(err)
		return err
	}

	r := bytes.NewReader(b)
	size := r.Size()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/gateway/v1/message/send_msg", serverAddr), r)
	if err != nil {
		logger.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.Header.Set("Authorization", token)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println(rsp.StatusCode)
	_ = rsp.Body.Close()
	v.Clear()
	v.SetCursor(v.Origin())
	g.Update(func(gg *gocui.Gui) error {
		v, err1 := gg.View("output")
		if err1 != nil {
			logger.Println("update err:", err1)
			return err1
		}
		fmt.Fprintln(v, "------")
		fmt.Fprintf(v, "Send|From:%v|Tp:%v|Content:%v\n", m.FromUser, m.ContentType, m.Content)
		return nil
	})
	return nil
}
