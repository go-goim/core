package main

import (
	"flag"
	"github.com/yusank/goim/apps/msg/internal/service"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/yusank/goim/apps/msg/internal/app"
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
		log.Fatal(err)
	}

	// register consumer
	c, err := mq.NewConsumer(&mq.ConsumerConfig{
		Addr:        application.ServerConfig.Mq.GetAddr(),
		Concurrence: 1,
		Subscriber:  service.GetMqMessageService(),
	})
	if err != nil {
		log.Fatal(err)
	}

	app.AddConsumer(c)
	if err = application.Run(); err != nil {
		log.Info(err)
	}

	application.Stop()
}
