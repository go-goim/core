package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-goim/core/pkg/waitgroup"
)

func main() {
	var (
		maxUser  int
		register bool
	)
	flag.IntVar(&maxUser, "max", 500, "max user")
	flag.BoolVar(&register, "register", false, "register user")
	flag.Parse()

	ip, err := getHostIP()
	if err != nil {
		panic(err)
	}

	ip += ":18071"
	addr = &ip

	users := make([]*user, maxUser)
	for i := 0; i < maxUser; i++ {
		u := &user{
			idx:   i,
			email: fmt.Sprintf("user%d@example.com", i),
		}

		if err := u.init(register); err != nil {
			panic(err)
		}

		users[i] = u
	}

	log.Println("user length:", len(users))
	wg := waitgroup.NewWaitGroup(20)
	for _, u := range users {
		u := u
		f := func() {
			u.handleFriendRequest()
		}
		wg.Add(f)
	}
	wg.Wait()
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
