package router

import (
	v1 "goAccounting/internal/api/v1"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes() {
	// 示例公开接口，可根据实际需求添加
	Public.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})

	userApi := new(v1.PublicApi)
	// 用户登录与注册
	Public.POST("/user/login", userApi.Login)
	Public.POST("/user/register", userApi.Register)
}
