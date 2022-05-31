package mid

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-goim/goim/pkg/log"
)

func Logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	c.Next()
	end := time.Now()
	latency := end.Sub(start)
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	log.Info("GinLog", "status", statusCode, "method", method, "path", path, "raw", raw,
		"latency", latency, "clientIP", clientIP)
}
