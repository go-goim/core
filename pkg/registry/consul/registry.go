package consul

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is consul registry option.
type Option func(*Registry)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// WithHeartbeat enable or disable heartbeat
func WithHeartbeat(enable bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.heartbeat = enable
		}
	}
}

// WithServiceResolver with endpoint function option.
func WithServiceResolver(fn ServiceResolver) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.resolver = fn
		}
	}
}

// WithHealthCheckInterval with healthcheck interval in seconds.
func WithHealthCheckInterval(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.healthcheckInterval = interval
		}
	}
}

// Config is consul registry config
type Config struct {
	*api.Config
}

// Registry is consul registry
type Registry struct {
	cli               *Client
	enableHealthCheck bool
}

// New creates consul registry
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		cli:               NewClient(apiClient),
		enableHealthCheck: true,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Register register service
func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

// Deregister deregister service
func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService return service by name
func (r *Registry) GetService(ctx context.Context, name string) (services []*registry.ServiceInstance, err error) {
	return r.cli.Service(ctx, name, true)
}

// ListServices return service list.
func (r *Registry) ListServices(ctx context.Context) (allServices []*registry.ServiceInstance, err error) {
	return r.cli.ListServices(ctx)

}

// Watch resolve service by name
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	return newConsulWatcher(ctx, r.cli, name)
}
