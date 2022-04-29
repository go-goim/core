package mid

import (
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	c.Next()
}
