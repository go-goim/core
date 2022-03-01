package main

import (
	"flag"

	"github.com/yusank/goim/apps/gateway/app"

	"github.com/go-kratos/kratos/v2/log"
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

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
