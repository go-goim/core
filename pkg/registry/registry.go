package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	clientv3 "go.etcd.io/etcd/client/v3"

	registryv1 "github.com/yusank/goim/api/config/registry/v1"
)

type RegisterDiscover interface {
	registry.Registrar
	registry.Discovery
}

func NewRegistry(regCfg *registryv1.Registry) (RegisterDiscover, error) {
	if cfg := regCfg.GetEtcd(); cfg != nil {
		return newEtcdRegistry(cfg)
	}

	if cfg := regCfg.GetConsul(); cfg != nil {
		return newConsulRegistry(cfg)
	}

	return nil, nil
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
