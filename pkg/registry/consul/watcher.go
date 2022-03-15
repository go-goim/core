package consul

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type consulWatcher struct {
	c           *Client
	wp          *watch.Plan
	next        chan []*registry.ServiceInstance
	serviceName string

	ctx    context.Context
	cancel context.CancelFunc
}

func newConsulWatcher(ctx context.Context, c *Client, name string) (registry.Watcher, error) {
	log.Println("watch called,name=", name)
	ctx2, cancel := context.WithCancel(ctx)
	cw := &consulWatcher{
		c:           c,
		ctx:         ctx2,
		cancel:      cancel,
		next:        make(chan []*registry.ServiceInstance, 10),
		serviceName: name,
	}

	wp, err := watch.Parse(map[string]interface{}{"type": "service", "service": name, "passingonly": true})
	if err != nil {
		return nil, err
	}

	cw.wp = wp
	wp.Handler = cw.serviceHandler
	go func() {
		if err1 := wp.RunWithClientAndHclog(c.cli, wp.Logger); err1 != nil {
			log.Println(err1)
		}
	}()

	return cw, nil
}

func (cw *consulWatcher) serviceHandler(_ uint64, data interface{}) {
	entries, ok := data.([]*api.ServiceEntry)
	if !ok {
		return
	}

	var finaleEntries = make([]*api.ServiceEntry, 0)
	for _, e := range entries {
		if e.Service.Service != cw.serviceName {
			continue
		}

		var del bool
		for _, check := range e.Checks {
			// delete the node if the status is critical
			if check.Status == "critical" {
				del = true
				break
			}
		}

		// if delete then skip the node
		if del {
			continue
		}

		finaleEntries = append(finaleEntries, e)
	}

	cw.next <- defaultResolver(context.TODO(), finaleEntries)
}

func (cw *consulWatcher) Next() ([]*registry.ServiceInstance, error) {
	select {
	case <-cw.ctx.Done():
		return nil, cw.ctx.Err()
	case r := <-cw.next:
		return r, nil
	}
}

func (cw *consulWatcher) Stop() error {
	cw.cancel()
	if cw.wp == nil {
		return nil
	}
	cw.wp.Stop()

	return nil
}
