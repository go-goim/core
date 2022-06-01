package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/apps/push/internal/app"
	"github.com/go-goim/core/apps/push/internal/router"
	"github.com/go-goim/core/apps/push/internal/service"
	"github.com/go-goim/core/pkg/cmd"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/mid"
)

var (
	jwtSecret string
	agentID   string // use hostname as agentID in default.
)

func init() {
	agentID, _ = os.Hostname()
	log.Debug("agent id", "agentID", agentID)
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
		log.Fatal("initApplication got err", "error", err)
	}

	// register grpc
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())

	// register router
	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	router.RegisterRouter(g.Group("/push/service"))
	application.HTTPSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Fatal("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Info("graceful shutdown got error", "error", err)
	}
}
