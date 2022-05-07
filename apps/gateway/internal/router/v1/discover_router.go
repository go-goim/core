package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/response"
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
		response.ErrorResp(c, fmt.Errorf("uid is empty"))
		return
	}

	agentID, err := service.LoadMatchedPushServer(context.Background())
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{
		"agentId": agentID,
	})
}
