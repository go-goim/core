package main

import (
	"context"

	"github.com/yusank/goim/apps/user/internal/app"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
)

func main() {
	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("InitApplication got err", "error", err)
	}

	if err = application.Run(); err != nil {
		log.Error("application run error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown error", "error", err)
	}
}
