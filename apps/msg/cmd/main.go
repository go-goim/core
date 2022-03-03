package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/app"
	"github.com/yusank/goim/apps/msg/service"
)

var (
	flagconf string
)

// todo: 整合合并重复的代码 包括初始化 app 以及 config 的解析

func init() {
	flag.StringVar(&flagconf, "conf", "../config", "config path, eg: --conf config.yaml")
}

func main() {
	flag.Parse()
	application, err := app.InitApplication(flagconf)
	if err != nil {
		log.Fatal(err)
	}

	messagev1.RegisterSendMeesagerServer(application.GrpcServer, &service.SendMessage{})

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
