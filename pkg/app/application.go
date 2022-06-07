package app

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"go.uber.org/atomic"
	ggrpc "google.golang.org/grpc"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	redisv8 "github.com/go-redis/redis/v8"

	"github.com/go-goim/core/pkg/cmd"
	"github.com/go-goim/core/pkg/config"
	"github.com/go-goim/core/pkg/db/mysql"
	"github.com/go-goim/core/pkg/db/redis"
	"github.com/go-goim/core/pkg/errors"
	"github.com/go-goim/core/pkg/mq"
	"github.com/go-goim/core/pkg/registry"
)

// Application is a common app entry.
// All apps can use this Application as a base entry and add own fields and methods
//  in their own app packages.
type Application struct {
	Core     *kratos.App
	Register registry.RegisterDiscover
	HTTPSrv  *http.Server
	GrpcSrv  *grpc.Server
	Config   *config.Config
	Producer mq.Producer
	Redis    *redisv8.Client
	Consumer []mq.Consumer

	host    string
	options *options
}

type options struct {
	metadata map[string]string
}

func newOptions(opts ...Option) *options {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	return opt
}

type Option func(o *options)

func WithMetadata(k, v string) Option {
	return func(o *options) {
		if o.metadata == nil {
			o.metadata = make(map[string]string)
		}
		o.metadata[k] = v
	}
}

var (
	useHostIP bool
)

func init() {
	cmd.GlobalFlagSet.BoolVar(&useHostIP, "use-host-ip", true, "use host ip")
}

var (
	initFlag atomic.Bool
)

func AssertApplication() {
	if !initFlag.Load() {
		panic("Application not initialized")
	}
}

func InitApplication(opts ...Option) (*Application, error) {
	// only can call this func once, if call twice will be panic
	if !initFlag.CAS(false, true) {
		panic("Application already initialized, don't call init function twice")
	}
	// init config
	cfg := config.InitConfig()

	a := &Application{
		Config:  cfg,
		options: newOptions(opts...),
	}

	if err := a.initHost(); err != nil {
		return nil, err
	}

	var servers = make([]transport.Server, 0)
	// init http server
	if err := a.initHTTPServer(); err != nil {
		return nil, err
	}
	if a.HTTPSrv != nil {
		servers = append(servers, a.HTTPSrv)
	}

	// init grpc server
	if err := a.initGrpcServer(); err != nil {
		return nil, err
	}
	if a.GrpcSrv != nil {
		servers = append(servers, a.GrpcSrv)
	}

	// init mq
	if err := a.initMq(); err != nil {
		return nil, err
	}

	// init db
	if err := a.initDB(); err != nil {
		return nil, err
	}

	// init kratos
	if err := a.initKratos(servers); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) initHost() error {
	if !useHostIP {
		a.host = ""
		return nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		// check the address type and do not use ipv6
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		// check the network type
		if ipnet.IP.IsLoopback() {
			continue
		}

		if ipnet.IP.To4() != nil {
			a.host = ipnet.IP.String()
			return nil
		}
	}

	return fmt.Errorf("not found host ip")
}

func (a *Application) initHTTPServer() error {
	if a.Config.SrvConfig.Http == nil {
		return nil
	}

	var timeout = time.Second
	if a.Config.SrvConfig.Http.GetTimeout() != nil && a.Config.SrvConfig.Http.GetTimeout().IsValid() {
		timeout = a.Config.SrvConfig.Http.GetTimeout().AsDuration()
	}

	portStr := strconv.Itoa(int(a.Config.SrvConfig.Http.GetPort()))
	httpSrv := http.NewServer(
		http.Address(net.JoinHostPort(a.host, portStr)),
		http.Middleware(
			recovery.Recovery(),
		),
		http.Timeout(timeout),
	)
	a.HTTPSrv = httpSrv

	return nil
}

func (a *Application) initGrpcServer() error {
	if a.Config.SrvConfig.Grpc == nil {
		return nil
	}

	var timeout = time.Second
	if a.Config.SrvConfig.Grpc.GetTimeout() != nil && a.Config.SrvConfig.Grpc.GetTimeout().IsValid() {
		timeout = a.Config.SrvConfig.Grpc.GetTimeout().AsDuration()
	}

	portStr := strconv.Itoa(int(a.Config.SrvConfig.Grpc.GetPort()))
	grpcSrv := grpc.NewServer(
		grpc.Address(net.JoinHostPort(a.host, portStr)),
		grpc.Middleware(
			recovery.Recovery(),
		),
		grpc.Timeout(timeout),
		grpc.Options(
			ggrpc.InitialWindowSize(1024*1024*1024),     // 1GB
			ggrpc.InitialConnWindowSize(1024*1024*1024), // 1GB
			ggrpc.MaxConcurrentStreams(1024),
		),
	)
	a.GrpcSrv = grpcSrv

	return nil
}

func (a *Application) initMq() error {
	if a.Config.SrvConfig.Mq == nil || len(a.Config.SrvConfig.Mq.GetAddr()) == 0 {
		return nil
	}

	p, err := mq.NewProducer(&mq.ProducerConfig{
		Retry: int(a.Config.SrvConfig.Mq.GetMaxRetry()),
		Addr:  a.Config.SrvConfig.Mq.GetAddr(),
	})
	if err != nil {
		return err
	}

	a.Producer = p

	return nil
}

func (a *Application) initDB() error {
	if err := a.initRedis(); err != nil {
		return err
	}

	if err := a.initMysql(); err != nil {
		return err
	}

	return nil
}

func (a *Application) initRedis() error {
	if a.Config.SrvConfig.GetRedis() == nil {
		return nil
	}

	rdb, err := redis.NewRedis(redis.WithConfig(a.Config.SrvConfig.GetRedis()))
	if err != nil {
		return err
	}

	a.Redis = rdb

	return nil
}

func (a *Application) initMysql() error {
	if a.Config.SrvConfig.GetMysql() == nil {
		return nil
	}

	err := mysql.InitDB(mysql.WithConfig(a.Config.SrvConfig.GetMysql()), mysql.Debug(a.Config.Debug()))
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initKratos(servers []transport.Server) error {
	a.initMetadata()

	var options = []kratos.Option{
		kratos.Name(a.Config.SrvConfig.GetName()),
		kratos.Version(a.Config.SrvConfig.GetVersion()),
		kratos.Server(
			servers...,
		),
		kratos.Metadata(
			a.options.metadata,
		),
	}

	reg, err := registry.NewRegistry(a.Config.RegConfig.Registry)
	if err != nil {
		return err
	}

	if reg != nil {
		a.Register = reg
		options = append(options, kratos.Registrar(reg))
	}

	a.Core = kratos.New(
		options...,
	)

	return nil
}

func (a *Application) initMetadata() {
	// metadata
	metadata := make(map[string]string)
	if a.Config.SrvConfig.GetMetadata() != nil {
		metadata = a.Config.SrvConfig.GetMetadata()
	}

	if len(a.options.metadata) > 0 {
		for k, v := range a.options.metadata {
			metadata[k] = v
		}
	}

	a.options.metadata = metadata
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
