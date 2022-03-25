package app

import (
	"github.com/yusank/goim/pkg/app"
	"github.com/yusank/goim/pkg/registry"
)

type Application struct {
	*app.Application
}

var (
	application *Application
)

func InitApplication(confPath string) (*Application, error) {
	cfg := app.ParseConfig(confPath)
	// do some own biz logic if needed
	a, err := app.InitApplication(cfg)
	if err != nil {
		return nil, err
	}

	application = &Application{Application: a}
	return application, nil
}

func GetApplication() *Application {
	return application
}

func GetRegister() registry.RegisterDiscover {
	return application.Register
}
