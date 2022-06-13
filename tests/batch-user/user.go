package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	friendpb "github.com/go-goim/api/user/friend/v1"
	userv1 "github.com/go-goim/api/user/v1"

	"github.com/go-goim/core/pkg/response"
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

	if data.Code != 0 {
		return fmt.Errorf("register user=%s, err= %v", u.email, data.Reason)
	}

	return nil
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

	user := new(userv1.User)
	if err := json.Unmarshal(b, user); err != nil {
		return err
	}

	u.uid = user.Uid
	return nil
}

func (u *user) addFriend(fuid string) error {
	str := fmt.Sprintf(`{"friend_uid":"%s"}`, fuid)

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
	str := fmt.Sprintf(`{"email":"%s"}`, email)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", *addr, queryFriend), strings.NewReader(str))
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
		return "", fmt.Errorf("query friend user=%s, friend=%s, err= %v", u.email, email, data.Reason)
	}

	b, err := json.Marshal(data.Data)
	if err != nil {
		return "", err
	}

	user := new(userv1.User)
	if err := json.Unmarshal(b, user); err != nil {
		return "", err
	}

	return user.Uid, nil
}

func (u *user) queryFriendRequests() (ids []int64, err error) {
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
		List []*friendpb.FriendRequest `json:"list"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	for _, v := range result.List {
		if v.Status == friendpb.FriendRequestStatus_REQUESTED {
			ids = append(ids, v.Id)
		}
	}

	log.Printf("query friend request user=%s, len(ids)=%v", u.email, len(ids))
	return ids, nil
}

func (u *user) acceptFriend(requestID int64) error {
	str := fmt.Sprintf(`{"uid":"%s", "friend_request_id":"%d", "action": 0}`, u.uid, requestID)

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
		return fmt.Errorf("accept friend user=%s, request_id=%d, err= %v", u.email, requestID, data.Reason)
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
	users := u.randomUsers(max)
	for _, email := range users {
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
