package app

import (
	"github.com/go-goim/core/pkg/app"
)

type Application struct {
	*app.Application
	agentID string
}

var (
	application *Application
)

func InitApplication(agentID string) (*Application, error) {
	// do some own biz logic if needed
	application = &Application{agentID: agentID}

	a, err := app.InitApplication(app.WithMetadata("agentID", agentID))
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
