package consul

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-goim/core/pkg/log"

	"github.com/hashicorp/consul/api"
)

// Client is consul client config
type Client struct {
	cli    *api.Client
	ctx    context.Context
	cancel context.CancelFunc

	// resolve service entry endpoints
	resolver ServiceResolver
	// healthcheck time interval in seconds
	healthcheckInterval int
	// heartbeat enable heartbeat
	heartbeat bool
}

// NewClient creates consul client
func NewClient(cli *api.Client) *Client {
	c := &Client{
		cli:                 cli,
		resolver:            defaultResolver,
		healthcheckInterval: 5,
		heartbeat:           true,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		var version string
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}

	return services
}

// ServiceResolver is used to resolve service endpoints
type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance

// Service get services from consul
func (c *Client) Service(ctx context.Context, service string, passingOnly bool) ([]*registry.ServiceInstance, error) {
	opts := &api.QueryOptions{
		WaitTime: time.Second * 30,
	}
	opts = opts.WithContext(ctx)
	entries, _, err := c.cli.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		log.Error("consul get service error", "err", err, "service", service)
		return nil, err
	}
	return c.resolver(ctx, entries), nil
}

// ListServices get services from consul
func (c *Client) ListServices(ctx context.Context) ([]*registry.ServiceInstance, error) {
	opts := &api.QueryOptions{
		WaitTime: time.Second * 30,
	}
	opts = opts.WithContext(ctx)
	rsp, _, err := c.cli.Catalog().Services(opts)
	if err != nil {
		return nil, err
	}

	var services = make([]*registry.ServiceInstance, 0, len(rsp))

	for service := range rsp {
		services = append(services, &registry.ServiceInstance{Name: service})
	}

	return services, nil
}

// Register register service instance to consul
func (c *Client) Register(_ context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress)
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)

		checkAddresses = append(checkAddresses, net.JoinHostPort(addr, strconv.FormatUint(port, 10)))
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}

	var healthCheckInterval = c.healthcheckInterval
	if healthCheckInterval == 0 {
		healthCheckInterval = 5
	}

	if enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            address,
				Interval:                       fmt.Sprintf("%ds", healthCheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", healthCheckInterval*10),
				Timeout:                        "5s",
			})
		}
	}
	if c.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID:                        "service:" + svc.ID,
			TTL:                            fmt.Sprintf("%ds", healthCheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", healthCheckInterval*10),
		})
	}

	err := c.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	if c.heartbeat {
		go func() {
			time.Sleep(time.Second)
			err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			if err != nil {
				log.Error("consul update ttl heartbeat to consul failed", "err", err)
			}
			ticker := time.NewTicker(time.Second * time.Duration(healthCheckInterval))
			defer ticker.Stop()
			for {
				select {
				case <-c.ctx.Done():
					return
				default:
				}

				// when select has multiple cases, it executes randomly.
				select {
				case <-c.ctx.Done():
					return
				case <-ticker.C:
					if c.ctx.Err() != nil {
						log.Info("consul heartbeat canceled")
						return
					}

					err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
					if err != nil {
						log.Error("consul update ttl heartbeat to consul failed", "err", err)
					}
				}
			}
		}()
	}
	return nil
}

// Deregister deregister service by service ID
func (c *Client) Deregister(_ context.Context, serviceID string) error {
	c.cancel()
	return c.cli.Agent().ServiceDeregister(serviceID)
}
