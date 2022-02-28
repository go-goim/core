package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/app/push/conf"
	"github.com/yusank/goim/app/push/service"
)

var (
	Name    = "goim.gateway.service"
	Version = "v0.0.0"
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../config", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	cfg := conf.ParseConfig(flagconf)
	s := &service.PushMessager{}

	var servers = make([]transport.Server, 0)
	if cfg.Http != nil {
		// debug and metrics
		servers = append(servers, http.NewServer(
			http.Address(":8000"),
			http.Middleware(
				recovery.Recovery(),
			),
		))
	}
	if cfg.Grpc != nil {
		// services
		grpcSrv := grpc.NewServer(
			grpc.Address(":9000"),
			grpc.Middleware(
				recovery.Recovery(),
			),
		)
		servers = append(servers, grpcSrv)
		messagev1.RegisterPushMessagerServer(grpcSrv, s)
	}

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			servers...,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
