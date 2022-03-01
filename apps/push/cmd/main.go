package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/push/conf"
	"github.com/yusank/goim/apps/push/router"
	"github.com/yusank/goim/apps/push/service"
	"github.com/yusank/goim/pkg/registry"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../config", "config path, eg: --conf config.yaml")
}

func main() {
	flag.Parse()
	cfg, regCfg := conf.ParseConfig(flagconf)
	s := &service.PushMessager{}

	var servers = make([]transport.Server, 0)
	if cfg.Http != nil {
		g := gin.New()
		router.RegisterRouter(g.Group("/gateway/api"))
		httpSrv := http.NewServer(
			http.Address(fmt.Sprintf("%s:%d", cfg.Http.GetAddr(), cfg.Http.GetPort())),
			http.Middleware(
				recovery.Recovery(),
			),
		)
		httpSrv.HandlePrefix("/", g)
		servers = append(servers, httpSrv)
	}
	if cfg.Grpc != nil {
		// services
		grpcSrv := grpc.NewServer(
			grpc.Address(fmt.Sprintf("%s:%d", cfg.Grpc.GetAddr(), cfg.Grpc.GetPort())),
			grpc.Middleware(
				recovery.Recovery(),
			),
		)
		servers = append(servers, grpcSrv)
		messagev1.RegisterPushMessagerServer(grpcSrv, s)
	}

	var options = []kratos.Option{
		kratos.Name(cfg.GetName()),
		kratos.Version(cfg.GetVersion()),
		kratos.Server(
			servers...,
		),
		kratos.Metadata(cfg.GetMetadata()),
	}

	reg, err := registry.NewRegistry(regCfg)
	if err != nil {
		panic(err)
	}

	if reg != nil {
		options = append(options, kratos.Registrar(reg))
	}

	app := kratos.New(
		options...,
	)

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
