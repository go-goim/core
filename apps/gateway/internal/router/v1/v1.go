package v1

import "github.com/gin-gonic/gin"

func Register(g *gin.RouterGroup) {
	g.GET("/discover", handleDiscoverPushServer)
	g.POST("/offline_msg/query", handleQueryOfflineMessage)
	// msg
	msg := g.Group("/msg")
	msg.POST("", handleSendSingleUserMsg)
	msg.POST("/broadcast", handleSendBroadcastMsg)
}
