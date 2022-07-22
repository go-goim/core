package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	friendpb "github.com/go-goim/api/user/friend/v1"
	"github.com/go-goim/core/pkg/web/response"
)

var (
	addr = flag.String("addr", "127.0.0.1:18071", "gateway addr")
)

const (
	registerURI           = "/gateway/v1/user/register"
	loginURI              = "/gateway/v1/user/login"
	queryFriend           = "/gateway/v1/user/query"
	addFriendURI          = "/gateway/v1/user/friend/add"
	acceptFriendURI       = "/gateway/v1/user/friend/accept"
	queryFriendRequestURI = "/gateway/v1/user/friend/request/list"
)

type user struct {
	idx   int
	uid   string
	email string
	token string
}

func (u *user) register() error {
	req := fmt.Sprintf(`{"email":"%s","password":"123456","name":"%s"}`, u.email, strings.TrimSuffix(u.email, "@example.com"))

	resp, err := http.Post(fmt.Sprintf("http://%s%s", *addr, registerURI), "application/json", strings.NewReader(req))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	if data.Code == 20002 {
		return nil
	}

	if data.Code != 0 {
		return fmt.Errorf("register user=%s, code= %v,reaso=%v", u.email, data.Code, data.Reason)
	}

	return nil
}

type User struct {
	UID         string  `json:"uid"` // 这里直接用 string 接,因为这个测试程序相当于客户端,客户端不需要解析
	Name        string  `json:"name" example:"user1"`
	Avatar      string  `json:"avatar" example:"https://www.example.com/avatar.png"`
	Email       *string `json:"email,omitempty" example:"abc@example.com"`
	Phone       *string `json:"phone,omitempty" example:"13800138000"`
	ConnectURL  *string `json:"connectUrl,omitempty" example:"ws://10.0.0.1:8080/ws"`
	LoginStatus int32   `json:"loginStatus" example:"0"`
}

func (u *user) login() error {
	req := fmt.Sprintf(`{"email":"%s","password":"123456"}`, u.email)

	resp, err := http.Post(fmt.Sprintf("http://%s%s", *addr, loginURI), "application/json", strings.NewReader(req))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	if data.Code != 0 {
		return fmt.Errorf("login user=%s, err= %v", u.email, data.Reason)
	}

	u.token = resp.Header.Get("Authorization")

	b, err := json.Marshal(data.Data)
	if err != nil {
		return err
	}

	user := new(User)
	if err := json.Unmarshal(b, user); err != nil {
		return err
	}

	u.uid = user.UID
	return nil
}

func (u *user) addFriend(fuid string) error {
	str := fmt.Sprintf(`{"friendUid":"%s"}`, fuid)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", *addr, addFriendURI), strings.NewReader(str))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	if data.Code != 0 {
		return fmt.Errorf("add friend user=%s, friend=%s, err= %v", u.email, fuid, data.Reason)
	}

	return nil
}

func (u *user) queryFriend(email string) (uid string, err error) {
	v := url.Values{}
	v.Set("email", email)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s?%s", *addr, queryFriend, v.Encode()), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return "", err
	}

	if data.Code != 0 {
		return "", fmt.Errorf("query friend user=%s, friend=%s, code= %v reason=%v", u.email, email, data.Code, data.Reason)
	}

	b, err := json.Marshal(data.Data)
	if err != nil {
		log.Println(string(b))
		panic(err)
		return "", err
	}

	user := new(User)
	if err := json.Unmarshal(b, user); err != nil {
		log.Println(string(b))
		panic(err)
		return "", err
	}

	return user.UID, nil
}

func (u *user) queryFriendRequests() (ids []uint64, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", *addr, queryFriendRequestURI), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	if data.Code != 0 {
		return nil, fmt.Errorf("query friend user=%s, err= %v", u.email, data.Reason)
	}

	b, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}

	var result struct {
		List []struct {
			ID     uint64 `json:"id"`
			Status int32  `json:"status"`
		} `json:"list"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println(string(b))
		panic(err)
		return nil, err
	}

	for _, v := range result.List {
		if v.Status == int32(friendpb.FriendRequestStatus_REQUESTED) {
			ids = append(ids, v.ID)
		}
	}

	return ids, nil
}

func (u *user) acceptFriend(requestID uint64) error {
	str := fmt.Sprintf(`{"friendRequestId":%d}`, requestID)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", *addr, acceptFriendURI), strings.NewReader(str))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data = new(response.Response)
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	if data.Code != 0 {
		return fmt.Errorf("accept friend user=%s, request_id=%d, err= %v", u.email, requestID, data.Message)
	}

	return nil
}

func (u *user) removeFriend(requestID int64) error {
	return nil
}

func (u *user) init(reg bool) error {
	if reg {
		if err := u.register(); err != nil {
			log.Printf("register user=%s, err= %s", u.email, err.Error())
			return err
		}
	}

	if err := u.login(); err != nil {
		log.Printf("login user=%s, err= %s", u.email, err.Error())
		return err
	}

	return nil
}

func (u *user) handleFriendRequest() {
	ids, err := u.queryFriendRequests()
	if err != nil {
		log.Printf("query friend request user=%s, err= %s", u.email, err.Error())
		return
	}

	for _, id := range ids {
		if err := u.acceptFriend(id); err != nil {
			log.Printf("accept friend request user=%s, err= %s", u.email, err.Error())
			return
		}
	}
}

func (u *user) run(max int) {
	log.Println("start", u.uid)
	users := u.randomUsers(max)
	for i, email := range users {
		if i%10 == 0 {
			log.Println("from", u.uid, "to", email, "idx", i)
		}
		uid, err := u.queryFriend(email)
		if err != nil {
			log.Printf("query friend user=%s, friend=%s, err= %s", u.email, email, err.Error())
			continue
		}

		if err := u.addFriend(uid); err != nil {
			log.Printf("add friend user=%s, friend=%s, err= %s", u.email, email, err.Error())
			continue
		}
	}
	log.Println("end", u.uid)
}

func (u *user) randomUsers(max int) []string {
	var maxFriend = 1000
	if max < 1000 {
		maxFriend = max - 1
	}
	users := make([]string, maxFriend)
	for i := 0; i < maxFriend; i++ {
		r := u.idx + i + 1
		if r >= max {
			r = r % max
		}
		if r == u.idx {
			continue
		}

		email := fmt.Sprintf("user%d@example.com", r)
		users[i] = email
	}

	return users
}
