package app

import (
	"github.com/go-goim/core/pkg/app"
)

type Application struct {
	*app.Application
}

var (
	application *Application
)

func InitApplication() (*Application, error) {
	// do some own biz logic if needed
	a, err := app.InitApplication()
	if err != nil {
		return nil, err
	}

	application = &Application{Application: a}
	return application, nil
}

func GetApplication() *Application {
	return application
}
