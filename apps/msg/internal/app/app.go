package app

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	redisv8 "github.com/go-redis/redis/v8"

	"github.com/yusank/goim/pkg/db/redis"
	"github.com/yusank/goim/pkg/mq"
	"github.com/yusank/goim/pkg/registry"
)

type Application struct {
	Core           *kratos.App
	Register       registry.RegisterDiscover
	ServerConfig   *Config
	RegisterConfig *Registry
	HttpServer     *http.Server
	GrpcServer     *grpc.Server
	Redis          *redisv8.Client
	Producer       mq.Producer
}

var (
	application *Application
	onceChan    = make(chan struct{}, 1)
)

// todo : app 层面尽可能不要引入当前 app 的包，否则很容易造成 import cycle 问题,最好是把 conf 和 app 包合并

func InitApplication(confPath string) (*Application, error) {
	// only can call this func once, if call twice will be panic
	onceChan <- struct{}{}
	close(onceChan)

	cfg, regCfg := ParseConfig(confPath)
	var servers = make([]transport.Server, 0)

	application = &Application{
		ServerConfig:   cfg,
		RegisterConfig: regCfg,
	}

	if cfg.Http != nil {
		httpSrv := http.NewServer(
			http.Address(fmt.Sprintf("%s:%d", cfg.Http.GetAddr(), cfg.Http.GetPort())),
			http.Middleware(
				recovery.Recovery(),
			),
		)
		application.HttpServer = httpSrv
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
		application.GrpcServer = grpcSrv
		servers = append(servers, grpcSrv)
	}

	var options = []kratos.Option{
		kratos.Name(cfg.GetName()),
		kratos.Version(cfg.GetVersion()),
		kratos.Server(
			servers...,
		),
		kratos.Metadata(cfg.GetMetadata()),
	}

	reg, err := registry.NewRegistry(regCfg.Registry)
	if err != nil {
		return nil, err
	}

	if reg != nil {
		application.Register = reg
		options = append(options, kratos.Registrar(reg))
	}

	rdb, err := redis.NewRedis(redis.WithConfig(cfg.GetRedis()))
	if err != nil {
		return nil, err
	}
	application.Redis = rdb

	application.Core = kratos.New(
		options...,
	)

	return application, nil
}

func (a *Application) Run() error {
	return a.Core.Run()
}

func GetApplication() *Application {
	return application
}

func GetRegister() registry.RegisterDiscover {
	return application.Register
}

func GetServerConfig() *Config {
	return application.ServerConfig
}
