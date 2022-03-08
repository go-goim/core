package registry

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	clientv3 "go.etcd.io/etcd/client/v3"

	registryv1 "github.com/yusank/goim/api/config/registry/v1"
)

var (
	registerInstance           RegisterDiscover
	ErrNotInitRegisterInstance = errors.New("registry instance not assigned")
	ErrUnknownRegisterInfo     = errors.New("register info unknown")
)

func GetRegisterDiscover() (RegisterDiscover, error) {
	if registerInstance == nil {
		return nil, ErrNotInitRegisterInstance
	}

	return registerInstance, nil
}

// Register the registration.
func Register(ctx context.Context, service *registry.ServiceInstance) error {
	if registerInstance == nil {
		return ErrNotInitRegisterInstance
	}

	return registerInstance.Register(ctx, service)
}

// Deregister the registration.
func Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	if registerInstance == nil {
		return ErrNotInitRegisterInstance
	}

	return registerInstance.Deregister(ctx, service)
}

// GetService return the service instances in memory according to the service name.
func GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	if registerInstance == nil {
		return nil, ErrNotInitRegisterInstance
	}

	return registerInstance.GetService(ctx, serviceName)
}

// Watch creates a watcher according to the service name.
func Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if registerInstance == nil {
		return nil, ErrNotInitRegisterInstance
	}

	return registerInstance.Watch(ctx, serviceName)
}

type RegisterDiscover interface {
	registry.Registrar
	registry.Discovery
}

func NewRegistry(regCfg *registryv1.Registry) (RegisterDiscover, error) {
	var (
		f   func(cfg *registryv1.RegistryInfo) (RegisterDiscover, error)
		cfg *registryv1.RegistryInfo
	)

	if c := regCfg.GetEtcd(); c != nil {
		f = newEtcdRegistry
		cfg = c
	}

	if c := regCfg.GetConsul(); c != nil {
		f = newConsulRegistry
		cfg = c
	}

	if f == nil {
		return nil, ErrUnknownRegisterInfo
	}

	rd, err := f(cfg)
	if err != nil {
		return nil, err
	}

	registerInstance = rd
	return registerInstance, nil
}

func newEtcdRegistry(cfg *registryv1.RegistryInfo) (RegisterDiscover, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.GetAddr(),
		DialTimeout:          cfg.GetDialTimeoutSec().AsDuration(),
		DialKeepAliveTime:    cfg.GetDialKeepAliveTimeSec().AsDuration(),
		DialKeepAliveTimeout: cfg.GetDialKeepAliveTimeoutSec().AsDuration(),
	})
	if err != nil {
		return nil, err
	}

	return etcd.New(cli), nil
}

func newConsulRegistry(cfg *registryv1.RegistryInfo) (RegisterDiscover, error) {
	cli, err := api.NewClient(&api.Config{
		Address: cfg.GetAddr()[0],
		Scheme:  cfg.GetScheme(),
	})

	if err != nil {
		return nil, err
	}

	return consul.New(cli, consul.WithHeartbeat(false), consul.WithHealthCheck(false)), nil
}
