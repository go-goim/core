package main

import (
	"context"

	"github.com/gin-gonic/gin"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/router"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
)

func main() {
	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("initApplication got err", "error", err)
	}

	log.Info("gateway start", "addr", application.Config.SrvConfig.Http.Addr, "version", application.Config.SrvConfig.Version)

	// register grpc
	messagev1.RegisterSendMessagerServer(application.GrpcSrv, &service.SendMessageService{})

	g := gin.Default()
	router.RegisterRouter(g.Group("/gateway/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Error("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown got error", "error", err)
	}
}
