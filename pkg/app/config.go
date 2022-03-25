package app

import (
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	registryv1 "github.com/yusank/goim/api/config/registry/v1"
	configv1 "github.com/yusank/goim/api/config/v1"
)

// Config contains service config.
// Use this as a basic config and add own fields in own app packages if needed.
type Config struct {
	SrvConfig *ServiceConfig
	RegConfig *RegistryConfig
}

func (c *Config) Validate() error {
	return nil
}

// ServiceConfig contains service config
type ServiceConfig struct {
	*configv1.Service `json:",inline"`
	FilePath          string
}

func NewConfig() *ServiceConfig {
	return &ServiceConfig{
		Service: new(configv1.Service),
	}
}

// RegistryConfig contains registry config
type RegistryConfig struct {
	*registryv1.Registry `json:",inline"`
	FilePath             string
}

func NewRegistry() *RegistryConfig {
	return &RegistryConfig{
		Registry: new(registryv1.Registry),
	}
}

func ParseConfig(fp string) *Config {
	c := config.New(
		config.WithSource(
			file.NewSource(fp),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	cfg := NewConfig()
	// Unmarshal the config to struct
	if err := c.Scan(cfg); err != nil {
		panic(err)
	}
	cfg.FilePath = fp
	log.Printf("%+v", cfg)

	reg := NewRegistry()
	if err := c.Scan(reg); err != nil {
		panic(err)
	}
	reg.FilePath = fp
	log.Printf("%+v", reg)
	reg.Name = cfg.GetName()

	return &Config{
		SrvConfig: cfg,
		RegConfig: reg,
	}
}
