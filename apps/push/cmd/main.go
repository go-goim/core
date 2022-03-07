package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/app"
	"github.com/yusank/goim/apps/push/service"
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
	messagev1.RegisterPushMessagerServer(application.GrpcServer, &service.PushMessager{})

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
