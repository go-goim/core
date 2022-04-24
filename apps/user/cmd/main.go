package main

import (
	"context"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/app"
	"github.com/yusank/goim/apps/user/internal/service"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
)

func main() {
	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("InitApplication got err", "error", err)
	}

	// TODO: add registered grpc services to metadata in service registry.
	userv1.RegisterUserServiceServer(application.GrpcSrv, service.GetUserService())

	if err = application.Run(); err != nil {
		log.Error("application run error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown error", "error", err)
	}
}
