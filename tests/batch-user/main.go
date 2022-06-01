package main

import (
	"flag"
	"fmt"

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

	wg := waitgroup.NewWaitGroup(50)
	for _, u := range users {
		u := u
		f := func() {
			u.run(maxUser)
		}
		wg.Add(f)
	}
	wg.Wait()
}
