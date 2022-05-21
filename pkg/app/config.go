package app

import (
	"strings"

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
	SrvConfig *ServiceConfig
	RegConfig *RegistryConfig
}

func (c *Config) Validate() error {
	return nil
}

// Debug returns true if service is running in debug mode.
func (c *Config) Debug() bool {
	return c.SrvConfig.Log.Level == configv1.Level_DEBUG
}

// ServiceConfig contains service config
type ServiceConfig struct {
	*configv1.Service `json:",inline"`
	FilePath          string
	SimpleName        string
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

var (
	confPath string
)

func init() {
	cmd.GlobalFlagSet.StringVar(&confPath, "conf", "./configs", "set config path")
}

func ParseConfig() *Config {
	c := config.New(
		config.WithSource(
			file.NewSource(confPath),
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

	// validate config
	if err := cfg.ValidateAll(); err != nil {
		panic(err)
	}

	cfg.FilePath = confPath
	slice := strings.Split(cfg.GetName(), ".")
	if len(slice) < 3 {
		log.Fatal("invalid service name=", cfg.GetName())
	}

	cfg.SimpleName = slice[1]
	log.Debug("config content", "config", cfg)

	reg := NewRegistry()
	if err := c.Scan(reg); err != nil {
		panic(err)
	}

	// validate config
	if err := reg.ValidateAll(); err != nil {
		panic(err)
	}

	reg.FilePath = confPath
	log.Debug("registry content", "registry", reg)
	reg.Name = cfg.GetName()

	setLogger(cfg.SimpleName, cfg.Log)
	return &Config{
		SrvConfig: cfg,
		RegConfig: reg,
	}
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
	))

	log.SetKratosLogger(log.NewZapLogger(
		log.Level(level),
		log.OutputPath(logPath),
		log.FilenamePrefix("kratos."),
		log.EnableConsole(logConf != nil && logConf.EnableConsole),
		log.CallerDepth(6),
	))
}
