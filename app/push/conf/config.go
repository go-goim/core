package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	configv1 "github.com/yusank/goim/api/config/v1"
)

type Config struct {
	*configv1.Service
	Filepath string
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

	cfg := new(Config)
	// Unmarshal the config to struct
	if err := c.Scan(cfg); err != nil {
		panic(err)
	}

	return cfg
}
