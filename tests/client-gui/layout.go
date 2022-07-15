package main

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"strconv"
//	"strings"
//
//	"github.com/jroimartin/gocui"
//
//	messagev1 "github.com/go-goim/api/message/v1"
//)
//
//var layoutDone bool
//
//func layout(g *gocui.Gui) error {
//	g.Highlight = true
//	g.Cursor = true
//	g.SelFgColor = gocui.ColorGreen
//
//	x0, y0, x1, y1 := friend.getCoordinates(g.Size())
//	if v, err := g.SetView(friendsView, x0, y0, x1, y1); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//		v.Highlight = true
//		v.Title = "Friends"
//		if err = initFriends(g, v); err != nil {
//			logger.Println("1", err)
//			return err
//		}
//
//		if _, err := setCurrentViewOnTop(g, friendsView); err != nil {
//			logger.Println("2", err)
//			return err
//		}
//	}
//
//	x0, y0, x1, y1 = outpu.getCoordinates(g.Size())
//	if v, err := g.SetView(outputView, x0, y0, x1, y1); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//		v.Title = fmt.Sprintf("Current: %s | %s ", toUid, toName)
//	}
//
//	x0, y0, x1, y1 = input.getCoordinates(g.Size())
//	if v, err := g.SetView(inputView, x0, y0, x1, y1); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//		v.Title = " " + uid + " "
//		v.Editable = true
//		v.Wrap = true
//	}
//
//	layoutDone = true
//	return nil
//}
//
//func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
//	if _, err := g.SetCurrentView(name); err != nil {
//		return nil, err
//	}
//	return g.SetViewOnTop(name)
//}
//
//func nextView(g *gocui.Gui, v *gocui.View) error {
//	logger.Println("next view", v.Name())
//	defer func() {
//		logger.Println("after next view", g.CurrentView().Name())
//	}()
//	if v == nil || v.Name() == friendsView {
//		_, err := setCurrentViewOnTop(g, inputView)
//		if err != nil {
//			logger.Println("set current view err:", err)
//			return err
//		}
//		g.Cursor = true
//		return nil
//	}
//
//	_, err := setCurrentViewOnTop(g, friendsView)
//	if err != nil {
//		logger.Println("set current view err:", err)
//		return err
//	}
//	g.Cursor = true
//	return nil
//}
//
//func arrowDown(g *gocui.Gui, v *gocui.View) error {
//	if v != nil {
//		cx, cy := v.Cursor()
//		// check if cursor is at the bottom
//		if err := v.SetCursor(cx, cy+1); err != nil {
//			ox, oy := v.Origin()
//			if err := v.SetOrigin(ox, oy+1); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func arrowUp(g *gocui.Gui, v *gocui.View) error {
//	if v != nil {
//		cx, cy := v.Cursor()
//		ox, oy := v.Origin()
//		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
//			if err := v.SetOrigin(ox, oy-1); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func selectFriend(g *gocui.Gui, v *gocui.View) error {
//	logger.Println("select friend")
//	_, cy := v.Cursor()
//	//_, oy := v.Origin()
//	line, err := v.Line(cy)
//	if err != nil {
//		logger.Println(err)
//		return err
//	}
//
//	toUid = line[:strings.Index(line, "|")]
//	toName = line[strings.Index(line, "|")+1:]
//	logger.Println("toUid:", toUid)
//	_, err = g.SetCurrentView(inputView)
//	if err != nil {
//		logger.Println("set current view err:", err)
//		return err
//	}
//
//	g.Update(func(gg *gocui.Gui) error {
//		v, err1 := gg.View(outputView)
//		if err1 != nil {
//			logger.Println("update err:", err1)
//			return err1
//		}
//
//		v.Title = fmt.Sprintf("Current: %s | %s ", toUid, toName)
//		return nil
//	})
//
//	return nil
//}
//
//func initFriends(g *gocui.Gui, v *gocui.View) error {
//	for i, friend := range friends {
//		if i == 0 {
//			toName = friend.FriendName
//			toUid = friend.FriendUid
//		}
//
//		fmt.Fprintln(v, fmt.Sprintf("%v|%v", friend.FriendUid, friend.FriendName))
//	}
//
//	return nil
//}
//
//func resetInput(g *gocui.Gui, v *gocui.View) error {
//	buf := &bytes.Buffer{}
//	io.Copy(buf, v)
//	// todo need load friend list then send msg
//	m := &messagev1.SendMessageReq{
//		From:        curUser.Uid,
//		To:          toUid,
//		ContentType: 1,
//		Content:     strings.TrimSuffix(buf.String(), "\n"),
//	}
//	b, err := json.Marshal(&m)
//	if err != nil {
//		logger.Println(err)
//		return err
//	}
//
//	r := bytes.NewReader(b)
//	size := r.Size()
//
//	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/gateway/v1/message/send_msg", serverAddr), r)
//	if err != nil {
//		logger.Println(err)
//		return err
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
//	req.Header.Set("Authorization", token)
//	rsp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		logger.Println(err)
//		return err
//	}
//	logger.Println(rsp.StatusCode)
//	_ = rsp.Body.Close()
//	v.Clear()
//	v.SetCursor(v.Origin())
//	g.Update(func(gg *gocui.Gui) error {
//		v, err1 := gg.View("output")
//		if err1 != nil {
//			logger.Println("update err:", err1)
//			return err1
//		}
//		fmt.Fprintln(v, "------")
//		fmt.Fprintf(v, "Send|From:%v|Tp:%v|Content:%v\n", m.From, m.ContentType, m.Content)
//		return nil
//	})
//	return nil
//}
