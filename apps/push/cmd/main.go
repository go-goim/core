package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"

	messagev1 "github.com/go-goim/goim/api/message/v1"
	"github.com/go-goim/goim/apps/push/internal/app"
	"github.com/go-goim/goim/apps/push/internal/router"
	"github.com/go-goim/goim/apps/push/internal/service"
	"github.com/go-goim/goim/pkg/cmd"
	"github.com/go-goim/goim/pkg/graceful"
	"github.com/go-goim/goim/pkg/mid"
)

var (
	jwtSecret string
	agentID   string // use hostname as agentID in default.
)

func init() {
	agentID, _ = os.Hostname()
	cmd.GlobalFlagSet.StringVar(&jwtSecret, "jwt-secret", "", "jwt secret")
	cmd.GlobalFlagSet.StringVar(&agentID, "agent-id", agentID, "agent id")
}

func main() {
	if err := cmd.ParseFlags(); err != nil {
		panic(err)
	}

	if jwtSecret == "" {
		panic("jwt secret is empty")
	}
	mid.SetJwtHmacSecret(jwtSecret)

	if agentID == "" {
		panic("agent id is empty")
	}

	application, err := app.InitApplication(agentID)
	if err != nil {
		log.Fatal(err)
	}

	// register grpc
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())

	// register router
	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	router.RegisterRouter(g.Group("/push/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Infof("graceful shutdown error: %s", err)
	}
}
