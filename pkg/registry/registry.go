package registry

import (
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	registryv1 "github.com/yusank/goim/api/config/registry/v1"
	"github.com/yusank/goim/app/push/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewRegistry(regCfg *conf.Registry) (registry.Registrar, error) {
	if cfg := regCfg.GetEtcd(); cfg != nil {
		return newEtcdRegistry(cfg)
	}

	if cfg := regCfg.GetConsul(); cfg != nil {
		return newConsulRegistry(cfg)
	}

	return nil, nil
}

func newEtcdRegistry(cfg *registryv1.RegistryInfo) (registry.Registrar, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.GetAddr(),
		DialTimeout:          time.Second * time.Duration(cfg.GetDialTimeoutSec()),
		DialKeepAliveTime:    time.Second * time.Duration(cfg.GetDialKeepAliveTimeSec()),
		DialKeepAliveTimeout: time.Second * time.Duration(cfg.GetDialKeepAliveTimeoutSec()),
	})
	if err != nil {
		return nil, err
	}

	return etcd.New(cli), nil
}

func newConsulRegistry(cfg *registryv1.RegistryInfo) (registry.Registrar, error) {
	cli, err := api.NewClient(&api.Config{
		Address: cfg.GetAddr()[0],
		Scheme:  cfg.GetScheme(),
	})

	if err != nil {
		return nil, err
	}
	return consul.New(cli, consul.WithHeartbeat(false), consul.WithHealthCheck(false)), nil
}
