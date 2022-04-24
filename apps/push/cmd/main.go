package main

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/yusank/goim/pkg/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/router"
	"github.com/yusank/goim/apps/push/internal/service"
	"github.com/yusank/goim/pkg/graceful"
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
		log.Fatal("InitApplication got err", "error", err)
	}

	// register grpc
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())

	// register router
	g := gin.Default()
	router.RegisterRouter(g.Group("/push/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Error("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown got error", "error", err)
	}
}
