package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/push/internal/app"
	"github.com/yusank/goim/apps/push/internal/router"
	"github.com/yusank/goim/apps/push/internal/service"
	"github.com/yusank/goim/pkg/cmd"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/mid"
)

var jwtSecret string

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
		log.Fatal(err)
	}

	// register grpc
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())

	// register router
	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	router.RegisterRouter(g.Group("/push/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Infof("graceful shutdown error: %s", err)
	}
}
