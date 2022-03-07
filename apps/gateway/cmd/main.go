package main

import (
	"flag"

	"github.com/yusank/goim/apps/gateway/app"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../config", "config path, eg: --conf config.yaml")
}

func main() {
	flag.Parse()
	application, err := app.InitApplication(flagconf)
	if err != nil {
		log.Fatal(err)
	}

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
