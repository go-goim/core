package v1

import "github.com/gin-gonic/gin"

func Register(g *gin.RouterGroup) {
	g.GET("/discover", handleDiscoverPushServer)
	g.POST("/send_msg", handleSendMsg)
}
