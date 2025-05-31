package router

import (
	"goAccounting/initialize"
	v1 "goAccounting/internal/api/v1"

	"github.com/gin-gonic/gin"
)

var Public *gin.RouterGroup

func RegisterPublicRoutes() {
	if Public == nil {
		// 使用 RouterPrefix
		Public = Engine.Group(initialize.Config.System.RouterPrefix + "/public")
	}

	// 示例公开接口，可根据实际需求添加
	Public.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})

	userApi := new(v1.PublicApi)
	// 用户登录与注册
	Public.POST("/user/login", userApi.Login)
	Public.POST("/user/register", userApi.Register)
}
