package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/router"
	"github.com/yusank/goim/apps/push/internal/service"
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
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())
	g := gin.Default()
	router.RegisterRouter(g.Group("/push/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
