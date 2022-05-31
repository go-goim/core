package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	registryv1 "github.com/yusank/goim/api/config/registry/v1"
	configv1 "github.com/yusank/goim/api/config/v1"
	"github.com/yusank/goim/pkg/cmd"
	"github.com/yusank/goim/pkg/log"
)

// Config contains service config.
// Use this as a basic config and add own fields in own app packages if needed.
type Config struct {
	SrvConfig          *ServiceConfig
	RegConfig          *RegistryConfig
	ConfigSource       config.Source
	EnableConfigCenter bool
}

// Debug returns true if service is running in debug mode.
func (c *Config) Debug() bool {
	return c.SrvConfig.Log.Level == configv1.Level_DEBUG
}

// ServiceConfig contains service config
type ServiceConfig struct {
	*configv1.Service `json:",inline"`
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

func NewRegistryConfig() *RegistryConfig {
	return &RegistryConfig{
		Registry: new(registryv1.Registry),
	}
}

var (
	confPath           string
	enableConfigCenter bool
)

func init() {
	cmd.GlobalFlagSet.StringVar(&confPath, "conf", "./configs", "set config path")
	cmd.GlobalFlagSet.BoolVar(&enableConfigCenter, "enable-config-center", true, "enable config center")
}

func InitConfig() *Config {
	c := config.New(
		config.WithSource(
			file.NewSource(confPath),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	reg := NewRegistryConfig()
	if err := c.Scan(reg); err != nil {
		panic(err)
	}

	// validate config
	if err := reg.ValidateAll(); err != nil {
		panic(err)
	}

	reg.FilePath = confPath
	log.Debug("registry content", "registry", reg)

	// init config center
	if enableConfigCenter {
		if reg.GetConfigCenter() == nil {
			panic("remote config must be set")
		}

		if err := reg.GetConfigCenter().Validate(); err != nil {
			panic(err)
		}

		source, err := NewSource(reg.Registry)
		if err != nil {
			panic(err)
		}

		cfg := &Config{
			RegConfig:          reg,
			ConfigSource:       source,
			EnableConfigCenter: enableConfigCenter,
		}

		if err := cfg.readFromConfigCenter(); err != nil {
			panic(err)
		}

		log.Debug("config content", "config", cfg)
		return cfg
	}

	// read all config from local files
	cfg := NewConfig()
	if err := c.Scan(cfg); err != nil {
		panic(err)
	}

	// validate config
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	log.Debug("service content", "service", cfg)

	return &Config{
		SrvConfig: cfg,
		RegConfig: reg,
	}
}

func (c *Config) readFromConfigCenter() error {
	cfg := config.New(config.WithSource(c.ConfigSource))
	if err := cfg.Load(); err != nil {
		return err
	}

	c.SrvConfig = NewConfig()
	if err := cfg.Scan(c.SrvConfig); err != nil {
		return err
	}

	// validate config
	if err := c.SrvConfig.Validate(); err != nil {
		return err
	}

	setLogger(c.SrvConfig.Name, c.SrvConfig.Log)
	return nil
}

func setLogger(serviceName string, logConf *configv1.Log) {
	var (
		logPath = "./logs/" + serviceName
		level   = configv1.Level_INFO
	)

	if logConf != nil && logConf.LogPath != nil && len(*logConf.LogPath) != 0 {
		logPath = *logConf.LogPath
	}

	if logConf != nil {
		level = logConf.Level
	}

	log.SetLogger(log.NewZapLogger(
		log.Level(level),
		log.OutputPath(logPath),
		log.FilenamePrefix("app."),
		log.EnableConsole(logConf != nil && logConf.EnableConsole),
		log.CallerDepth(2),
		log.Meta("service", serviceName),
		log.Meta("source", "app"),
	))

	log.SetKratosLogger(log.NewZapLogger(
		log.Level(level),
		log.OutputPath(logPath),
		log.FilenamePrefix("kratos."),
		log.EnableConsole(logConf != nil && logConf.EnableConsole),
		log.CallerDepth(6),
		log.Meta("service", serviceName),
		log.Meta("source", "kratos"),
	))
}