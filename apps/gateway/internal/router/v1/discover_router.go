package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-goim/goim/apps/gateway/internal/service"
	"github.com/go-goim/goim/pkg/mid"
	"github.com/go-goim/goim/pkg/response"
	"github.com/go-goim/goim/pkg/router"
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

// @Summary 获取推送服务器
// @Description 获取推送服务器 agentID
// @Tags [gateway]discover
// @Produce  json
// @Param   token query string true "token"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Failure 401 {null} null
// @Router /gateway/v1/discovery/discover [get]
func (r *DiscoverRouter) handleDiscoverPushServer(c *gin.Context) {
	uid := c.GetString("uid")
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
		"agentID": agentID,
	})
}
