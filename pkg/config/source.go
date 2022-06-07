package config

import (
	"fmt"

	"github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
	clientv3 "go.etcd.io/etcd/client/v3"

	registryv1 "github.com/go-goim/api/config/registry/v1"

	"github.com/go-goim/core/pkg/config/consul"
)

// NewSource create a config source according to the registry info.
func NewSource(reg *registryv1.Registry) (s config.Source, err error) {
	if reg.GetEtcd() != nil {
		return newEtcdSource(reg)
	}

	if reg.GetConsul() != nil {
		return newConsulSource(reg)
	}

	return nil, fmt.Errorf("unknown registry info")
}

func newEtcdSource(reg *registryv1.Registry) (s config.Source, err error) {
	cfg := reg.GetEtcd()
	cc := reg.GetConfigCenter()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.GetAddr(),
		DialTimeout:          cfg.GetDialTimeoutSec().AsDuration(),
		DialKeepAliveTime:    cfg.GetDialKeepAliveTimeSec().AsDuration(),
		DialKeepAliveTimeout: cfg.GetDialKeepAliveTimeoutSec().AsDuration(),
	})
	if err != nil {
		return nil, err
	}

	// TODO: support set keys.
	return etcd.New(cli, etcd.WithPath(cc.GetPathPrefix()), etcd.WithPrefix(true))
}

func newConsulSource(reg *registryv1.Registry) (s config.Source, err error) {
	cfg := reg.GetConsul()
	cc := reg.GetConfigCenter()
	cli, err := api.NewClient(&api.Config{
		Address:    cfg.GetAddr()[0],
		Scheme:     cfg.GetScheme(),
		Datacenter: "dc1",
	})

	if err != nil {
		return nil, err
	}

	return consul.New(cli, consul.WithPathPrefix(cc.GetPathPrefix()),
		consul.WithPaths(cc.GetPaths()...),
		consul.WithFormat(cc.GetFormat()))
}
