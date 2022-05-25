package app

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/atomic"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	redisv8 "github.com/go-redis/redis/v8"

	"github.com/yusank/goim/pkg/db/mysql"
	"github.com/yusank/goim/pkg/db/redis"
	"github.com/yusank/goim/pkg/errors"
	"github.com/yusank/goim/pkg/mq"
	"github.com/yusank/goim/pkg/registry"
)

// Application is a common app entry.
// All apps can use this Application as a base entry and add own fields and methods
//  in their own app packages.
type Application struct {
	Core     *kratos.App
	Register registry.RegisterDiscover
	HTTPSrv  *http.Server
	GrpcSrv  *grpc.Server
	Config   *Config
	Producer mq.Producer
	Redis    *redisv8.Client
	Consumer []mq.Consumer
}

var (
	initFlag atomic.Bool
)

func AssertApplication() {
	if !initFlag.Load() {
		panic("Application not initialized")
	}
}

func InitApplication(cfg *Config) (*Application, error) {
	// only can call this func once, if call twice will be panic
	if !initFlag.CAS(false, true) {
		panic("Application already initialized, don't call init function twice")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	application := &Application{
		Config: cfg,
	}

	var servers = make([]transport.Server, 0)
	if cfg.SrvConfig.Http != nil {
		var timeout = time.Second
		if cfg.SrvConfig.Http.GetTimeout() != nil && cfg.SrvConfig.Http.GetTimeout().IsValid() {
			timeout = cfg.SrvConfig.Http.GetTimeout().AsDuration()
		}

		httpSrv := http.NewServer(
			http.Address(fmt.Sprintf("%s:%d", cfg.SrvConfig.Http.GetAddr(), cfg.SrvConfig.Http.GetPort())),
			http.Middleware(
				recovery.Recovery(),
			),
			http.Timeout(timeout),
		)
		application.HTTPSrv = httpSrv
		servers = append(servers, httpSrv)
	}
	if cfg.SrvConfig.Grpc != nil {
		var timeout = time.Second
		if cfg.SrvConfig.Grpc.GetTimeout() != nil && cfg.SrvConfig.Grpc.GetTimeout().IsValid() {
			timeout = cfg.SrvConfig.Grpc.GetTimeout().AsDuration()
		}
		// services
		grpcSrv := grpc.NewServer(
			grpc.Address(fmt.Sprintf("%s:%d", cfg.SrvConfig.Grpc.GetAddr(), cfg.SrvConfig.Grpc.GetPort())),
			grpc.Middleware(
				recovery.Recovery(),
			),
			grpc.Timeout(timeout),
		)
		application.GrpcSrv = grpcSrv
		servers = append(servers, grpcSrv)
	}

	if cfg.SrvConfig.Mq != nil && len(cfg.SrvConfig.Mq.GetAddr()) > 0 {
		p, err := mq.NewProducer(&mq.ProducerConfig{
			Retry: int(cfg.SrvConfig.Mq.GetMaxRetry()),
			Addr:  cfg.SrvConfig.Mq.GetAddr(),
		})
		if err != nil {
			return nil, err
		}

		application.Producer = p
	}

	if cfg.SrvConfig.GetRedis() != nil {
		rdb, err := redis.NewRedis(redis.WithConfig(cfg.SrvConfig.GetRedis()))
		if err != nil {
			return nil, err
		}

		application.Redis = rdb
	}

	if cfg.SrvConfig.GetMysql() != nil {
		err := mysql.InitDB(mysql.WithConfig(cfg.SrvConfig.GetMysql()), mysql.Debug(cfg.Debug()))
		if err != nil {
			return nil, err
		}
	}

	var options = []kratos.Option{
		kratos.Name(cfg.SrvConfig.GetName()),
		kratos.Version(cfg.SrvConfig.GetVersion()),
		kratos.Server(
			servers...,
		),
		kratos.Metadata(cfg.SrvConfig.GetMetadata()),
	}

	reg, err := registry.NewRegistry(cfg.RegConfig.Registry)
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
	if a.Producer != nil {
		if err := a.Producer.Start(); err != nil {
			return err
		}
	}
	for _, consumer := range a.Consumer {
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	return a.Core.Run()
}

func (a *Application) Shutdown(ctx context.Context) error {
	var (
		es                 = make(errors.ErrorSet, 0)
		checkCtxAndExecute = func(f func() error) {
			select {
			case <-ctx.Done():
				es = append(es, ctx.Err())
				return
			default:
			}

			if err := f(); err != nil {
				es = append(es, err)
			}
		}
	)

	if a.Producer != nil {
		checkCtxAndExecute(func() error {
			if err := a.Producer.Shutdown(); err != nil {
				return fmt.Errorf("shutdown producer error: %w", err)
			}

			return nil
		})
	}

	for _, consumer := range a.Consumer {
		consumer := consumer
		checkCtxAndExecute(func() error {
			if err := consumer.Shutdown(); err != nil {
				return fmt.Errorf("shutdown consumer error: %w", err)
			}

			return nil
		})
	}

	if a.Redis != nil {
		checkCtxAndExecute(func() error {
			if err := a.Redis.Close(); err != nil {
				return fmt.Errorf("close redis error: %w", err)
			}
			return nil
		})
	}

	return es.Err()
}

func (a *Application) AddConsumer(c mq.Consumer) {
	if a.Consumer == nil {
		a.Consumer = make([]mq.Consumer, 0)
	}

	a.Consumer = append(a.Consumer, c)
}
