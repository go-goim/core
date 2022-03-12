package app

import (
	"fmt"

	"github.com/yusank/goim/pkg/mq"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/yusank/goim/pkg/registry"
)

type Application struct {
	Core           *kratos.App
	Register       registry.RegisterDiscover
	HTTPSrv        *http.Server
	ServerConfig   *Config
	RegisterConfig *Registry
	Producer       mq.Producer
}

var (
	application *Application
	onceChan    = make(chan struct{}, 1)
)

func InitApplication(confPath string) (*Application, error) {
	// only can call this func once, if call twice will be panic
	onceChan <- struct{}{}
	close(onceChan)

	cfg, regCfg := ParseConfig(confPath)
	application = &Application{
		ServerConfig:   cfg,
		RegisterConfig: regCfg,
	}

	var servers = make([]transport.Server, 0)
	if cfg.Http != nil {
		httpSrv := http.NewServer(
			http.Address(fmt.Sprintf("%s:%d", cfg.Http.GetAddr(), cfg.Http.GetPort())),
			http.Middleware(
				recovery.Recovery(),
			),
		)
		application.HTTPSrv = httpSrv
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
	}

	mqCfg := &mq.ProducerConfig{
		Retry: int(cfg.Mq.GetMaxRetry()),
		Addr:  cfg.Mq.GetAddr(),
	}
	p, err := mq.NewProducer(mqCfg)
	if err != nil {
		return nil, err
	}

	application.Producer = p

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

	application.Core = kratos.New(
		options...,
	)

	return application, nil
}

func (a *Application) Run() error {
	if err := a.Producer.Start(); err != nil {
		return err
	}
	return a.Core.Run()
}

func (a *Application) Stop() {
	_ = a.Producer.Shutdown()
}

func GetRegister() registry.RegisterDiscover {
	return application.Register
}

func GetApplication() *Application {
	select {
	case <-onceChan:
		panic("application not init")
	default:
	}

	return application
}
