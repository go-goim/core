package mid

import (
	"context"

	"github.com/gin-gonic/gin"
)

func GetContext(c *gin.Context) context.Context {
	return context.Background()
}
