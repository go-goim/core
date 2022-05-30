package consul

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx        context.Context
	pathPrefix string
	paths      map[string]bool
	format     string
}

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithPathPrefix with config pathPrefix
func WithPathPrefix(p string) Option {
	return func(o *options) {
		o.pathPrefix = p
	}
}

// WithPaths with config paths
func WithPaths(paths ...string) Option {
	return func(o *options) {
		o.paths = make(map[string]bool)
		for _, key := range paths {
			o.paths[key] = true
		}
	}
}

func WithFormat(format string) Option {
	return func(o *options) {
		o.format = format
	}
}

type source struct {
	client  *api.Client
	options *options
}

func New(client *api.Client, opts ...Option) (config.Source, error) {
	options := &options{
		ctx:        context.Background(),
		pathPrefix: "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.pathPrefix == "" {
		return nil, errors.New("pathPrefix invalid")
	}

	return &source{
		client:  client,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	kv, _, err := s.client.KV().List(s.options.pathPrefix, nil)
	if err != nil {
		return nil, err
	}

	pathPrefix := s.options.pathPrefix
	if !strings.HasSuffix(s.options.pathPrefix, "/") {
		pathPrefix += "/"
	}
	kvs := make([]*config.KeyValue, 0)
	for _, item := range kv {
		k := strings.TrimPrefix(item.Key, pathPrefix)
		if k == "" {
			continue
		}

		if len(s.options.paths) > 0 {
			if _, ok := s.options.paths[k]; !ok {
				continue
			}
		}

		format := s.options.format
		if format == "" {
			format = strings.TrimPrefix(filepath.Ext(k), ".")
		}

		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: format,
		})
	}

	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s)
}
