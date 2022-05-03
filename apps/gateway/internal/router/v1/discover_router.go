package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/router"
)

type DiscoverRouter struct {
	router.Router
}

func NewDiscoverRouter() *DiscoverRouter {
	return &DiscoverRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *DiscoverRouter) Load(g *gin.RouterGroup) {
	g.GET("/discover", mid.AuthJwtCookie, r.handleDiscoverPushServer)
}

func (r *DiscoverRouter) handleDiscoverPushServer(c *gin.Context) {
	uid := c.GetHeader("uid")
	if uid == "" {
		log.Println("uid not found")
		c.JSON(http.StatusOK, gin.H{"err": "uid not found"})
		return
	}

	agentID, err := service.LoadMatchedPushServer(context.Background())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agentId": agentID})
}
