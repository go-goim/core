package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yusank/goim/apps/gateway/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/router"
	"github.com/yusank/goim/pkg/util"
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
		util.ErrorResp(c, fmt.Errorf("uid is empty"))
		return
	}

	agentID, err := service.LoadMatchedPushServer(context.Background())
	if err != nil {
		util.ErrorResp(c, err)
		return
	}

	util.SuccessResp(c, gin.H{
		"agentId": agentID,
	})
}
