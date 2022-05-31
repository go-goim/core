package main

import (
	"context"

	"github.com/gin-gonic/gin"

	messagev1 "github.com/go-goim/goim/api/message/v1"
	"github.com/go-goim/goim/apps/gateway/internal/app"
	"github.com/go-goim/goim/apps/gateway/internal/router"
	"github.com/go-goim/goim/apps/gateway/internal/service"
	"github.com/go-goim/goim/pkg/cmd"
	"github.com/go-goim/goim/pkg/graceful"
	"github.com/go-goim/goim/pkg/log"
	"github.com/go-goim/goim/pkg/mid"
)

var (
	jwtSecret = ""
)

func init() {
	cmd.GlobalFlagSet.StringVar(&jwtSecret, "jwt-secret", "", "jwt secret")
}

func main() {
	if err := cmd.ParseFlags(); err != nil {
		panic(err)
	}

	if jwtSecret == "" {
		panic("jwt secret is empty")
	}
	mid.SetJwtHmacSecret(jwtSecret)

	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("initApplication got err", "error", err)
	}

	log.Info("gateway start", "addr", application.Config.SrvConfig.Http.Addr, "version", application.Config.SrvConfig.Version)

	// register grpc
	messagev1.RegisterSendMessagerServer(application.GrpcSrv, &service.SendMessageService{})

	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
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
