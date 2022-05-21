package main

import (
	"context"

	"github.com/gin-gonic/gin"

	friendpb "github.com/yusank/goim/api/user/friend/v1"
	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/app"
	"github.com/yusank/goim/apps/user/internal/router"
	"github.com/yusank/goim/apps/user/internal/service"
	"github.com/yusank/goim/pkg/cmd"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/mid"
)

func main() {
	if err := cmd.ParseFlags(); err != nil {
		panic(err)
	}

	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("InitApplication got err", "error", err)
	}

	// TODO: add registered grpc services to metadata in service registry.
	userv1.RegisterUserServiceServer(application.GrpcSrv, service.GetUserService())
	friendpb.RegisterFriendServiceServer(application.GrpcSrv, service.GetFriendService())

	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	router.RegisterRouter(g.Group("/user/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Error("application run error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown error", "error", err)
	}
}
