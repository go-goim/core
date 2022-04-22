package main

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/router"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/graceful"

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

	log.Infow("gateway start", "addr", application.Config.SrvConfig.Http.Addr, "version", application.Config.SrvConfig.Version)

	// register grpc
	messagev1.RegisterSendMessagerServer(application.GrpcSrv, &service.SendMessageService{})

	g := gin.Default()
	router.RegisterRouter(g.Group("/gateway/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Error("application run error: %v", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Infof("graceful shutdown error: %v", err)
	}
}
