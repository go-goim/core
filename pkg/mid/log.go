package mid

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-goim/core/pkg/log"
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
	kvs := []interface{}{
		"tp", "HttpRequest",
		"status", statusCode,
		"method", method,
		"raw", raw,
		"latency", latency.Microseconds(),
		"clientIP", clientIP,
	}

	if len(c.Errors) > 0 {
		kvs = append(kvs, "error", c.Errors.ByType(gin.ErrorTypePrivate).String())
	}

	if c.GetString("uid") != "" {
		kvs = append(kvs, "uid", c.GetString("uid"))
	}

	log.Info(path, kvs...)
}
