package mid

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yusank/goim/pkg/log"
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
	log.Info(path, "tp", "HttpRequest", "status", statusCode, "method", method, "raw", raw,
		"latency", latency.Microseconds(), "clientIP", clientIP)
}
