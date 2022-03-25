package app

import (
	"github.com/yusank/goim/pkg/app"
	"github.com/yusank/goim/pkg/registry"
)

type Application struct {
	*app.Application
	// add own fields here
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

func GetRegister() registry.RegisterDiscover {
	return application.Register
}

func GetApplication() *Application {
	app.AssertApplication()
	return application
}
