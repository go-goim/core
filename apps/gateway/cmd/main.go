package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/router"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../config", "config path, eg: --conf config.yaml")
}

func main() {
	flag.Parse()
	application, err := app.InitApplication(flagconf)
	if err != nil {
		log.Fatal(err)
	}

	g := gin.New()
	router.RegisterRouter(g.Group("/gateway/service"))
	application.HttpSrv.HandlePrefix("/", g)

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}
