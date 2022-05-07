package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/yusank/goim/apps/push/internal/service"
	"github.com/yusank/goim/pkg/mid"
	"github.com/yusank/goim/pkg/resp"
	"github.com/yusank/goim/pkg/router"
)

type ConnRouter struct {
	router.Router
	upgrader websocket.Upgrader
}

func NewConnRouter() *ConnRouter {
	return &ConnRouter{
		Router: &router.BaseRouter{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (r *ConnRouter) Load(g *gin.RouterGroup) {
	g.GET("/ws", mid.AuthJwtCookie, r.wsHandler)
}

func (r *ConnRouter) wsHandler(c *gin.Context) {
	// todo use check uid/token middleware before this handler
	uid := c.GetHeader("uid")
	if uid == "" {
		resp.ErrorRespWithStatus(c, http.StatusUnauthorized, fmt.Errorf("invalid uid"))
		return
	}

	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		resp.ErrorRespWithStatus(c, http.StatusBadRequest, err)
		return
	}

	service.HandleWsConn(mid.GetContext(c), conn, uid)
}
