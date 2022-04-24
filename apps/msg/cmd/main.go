package main

import (
	"context"
	"flag"

	"github.com/yusank/goim/pkg/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
	"github.com/yusank/goim/apps/msg/internal/service"
	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/mq"
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
	messagev1.RegisterOfflineMessageServer(application.GrpcSrv, &service.OfflineMessageService{})

	// register consumer
	c, err := mq.NewConsumer(&mq.ConsumerConfig{
		Addr:        application.Config.SrvConfig.Mq.GetAddr(),
		Concurrence: 1,
		Subscriber:  service.GetMqMessageService(),
	})
	if err != nil {
		log.Fatal("NewConsumer got err", "error", err)
	}
	application.AddConsumer(c)

	if err = application.Run(); err != nil {
		log.Error("application run error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown error", "error", err)
	}
}
