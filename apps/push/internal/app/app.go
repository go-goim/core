package app

import (
	"github.com/google/uuid"

	"github.com/yusank/goim/pkg/app"
)

type Application struct {
	*app.Application
	agentID string
}

var (
	application *Application
)

func InitApplication() (*Application, error) {
	cfg := app.ParseConfig()
	// do some own biz logic if needed
	application = &Application{agentID: uuid.NewString()}

	cfg.SrvConfig.GetMetadata()["agentId"] = application.agentID
	a, err := app.InitApplication(cfg)
	if err != nil {
		return nil, err
	}
	application.Application = a

	return application, nil
}

func GetApplication() *Application {
	return application
}

func GetAgentID() string {
	return application.agentID
}
