package main

import (
	"context"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/router"
	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/cmd"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/mid"

	_ "github.com/swaggo/swag"

	_ "github.com/yusank/goim/swagger"
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
	router.RegisterRouter(g.Group("/gateway"))
	application.HTTPSrv.HandlePrefix("/", g)
	// register swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err = application.Run(); err != nil {
		log.Error("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown got error", "error", err)
	}
}
