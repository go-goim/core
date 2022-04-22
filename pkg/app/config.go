package app

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	registryv1 "github.com/yusank/goim/api/config/registry/v1"
	configv1 "github.com/yusank/goim/api/config/v1"
	"github.com/yusank/goim/pkg/log"
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
	log.Debug("config content", "config", cfg)

	reg := NewRegistry()
	if err := c.Scan(reg); err != nil {
		panic(err)
	}
	reg.FilePath = fp
	log.Debug("registry content", "registry", reg)
	reg.Name = cfg.GetName()

	setLogger(cfg.Log)
	return &Config{
		SrvConfig: cfg,
		RegConfig: reg,
	}
}

func setLogger(logConf *configv1.Log) {
	var (
		logPath = "./logs"
	)
	if logConf != nil && logConf.LogPath != nil {
		logPath = *logConf.LogPath
	}

	log.SetLogger(log.NewZapLogger(log.WithLevel(logConf.Level), log.WithOutputPath(logPath)))
}
