package main

//
//import (
//	"bytes"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"log"
//	"math/rand"
//	"net/http"
//	"os"
//	"os/signal"
//	"strconv"
//	"syscall"
//	"time"
//
//	messagev1 "github.com/go-goim/api/message/v1"
//)
//
//var (
//	gatewayAddr string
//	clientCount int
//)
//
//func init() {
//	flag.StringVar(&gatewayAddr, "s", "127.0.0.1:18071", "gateway server addr")
//	flag.IntVar(&clientCount, "c", 10, "client count")
//	flag.Parse()
//}
//
//func main() {
//	for i := 0; i < clientCount; i++ {
//		go runClient(fmt.Sprintf("user_%d", i))
//	}
//
//	var sigChan = make(chan os.Signal, 1)
//	signal.Notify(sigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGINT)
//	<-sigChan
//}
//
//func runClient(uid int64) {
//	err := initAndStoreClient(uid)
//	if err != nil {
//		panic(fmt.Errorf("uid=%s,err=%v", uid, err))
//	}
//	time.Sleep(time.Second)
//
//	var ticker = time.NewTicker(time.Millisecond * 50)
//	for range ticker.C {
//		if err = messageToRandomClient(uid); err != nil {
//			log.Println("messageToRandomClient got err=", err)
//		}
//	}
//}
//
//func messageToRandomClient(uid int64) error {
//	msg := messagev1.SendMessageReq{
//		From: uid,
//		//To:          randomUID(),
//		ContentType: messagev1.MessageContentType_Text,
//		Content:     randomMsg(20),
//	}
//
//	b, err := json.Marshal(&msg)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	r := bytes.NewReader(b)
//	size := r.Size()
//
//	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/gateway/service/v1/send_msg", gatewayAddr), r)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
//	rsp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	_ = rsp.Body.Close()
//
//	if rsp.StatusCode != http.StatusOK {
//		return fmt.Errorf("request send msg api got error status code=%v", rsp.StatusCode)
//	}
//
//	return nil
//}
//
//func randomUID() string {
//	return fmt.Sprintf("user_%d", rand.Intn(clientCount))
//}
//
//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//
//func randomMsg(n int) string {
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}
