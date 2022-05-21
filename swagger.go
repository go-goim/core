//go:generate swag init -g swagger.go -o ./swagger
package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/swaggo/swag"

	_ "github.com/yusank/goim/swagger"
)

// @title GoIM Swagger
// @version 1.0
// @description GoIM 服务器接口文档
// @termsOfService http://yusank.github.io/goim/

// @contact.name Yusank
// @contact.url https://yusank.space
// @contact.email yusankurban@gmail.com

// @license.name MIT
// @license.url https://github.com/yusank/goim/blob/main/LICENSE

// @BasePath /
func main() {
	g := gin.New()
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Run(":8080")
}
