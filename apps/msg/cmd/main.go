package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
	"github.com/yusank/goim/apps/msg/internal/service"
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
	messagev1.RegisterSendMeesagerServer(application.GrpcServer, &service.SendMessageService{})

	// register consumer
	c, err := mq.NewConsumer(&mq.ConsumerConfig{
		Addr:  application.ServerConfig.Mq.GetAddr(),
		Topic: "",
		Group: "",
		// should define interface which has methods to get topic, group and handle func
		Handler:     service.GetMqMessageService().HandleMqMessage,
		Concurrence: 1,
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
