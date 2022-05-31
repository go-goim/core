package app

import (
	"github.com/go-goim/goim/pkg/app"
)

type Application struct {
	*app.Application
	agentID string
}

var (
	application *Application
)

func InitApplication(agentID string) (*Application, error) {
	cfg := app.ParseConfig()
	// do some own biz logic if needed
	application = &Application{agentID: agentID}

	cfg.SrvConfig.GetMetadata()["agentID"] = application.agentID
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
