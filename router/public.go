package router

import (
	"goAccounting/global"
	"goAccounting/initialize"
	v1 "goAccounting/internal/api/v1"

	"github.com/gin-gonic/gin"
)

var Public *gin.RouterGroup

func RegisterPublicRoutes() {
	if Public == nil {
		Public = Engine.Group(initialize.Config.System.RouterPrefix + "/public")
	}

	// 健康检查端点，前端可以用来测试连接
	Public.GET("/health", func(ctx *gin.Context) {
		status := global.HealthCheck()
		ctx.JSON(200, gin.H{
			"message":  "Backend service is running",
			"status":   "ok",
			"services": status,
		})
	})

	// 示例公开接口，可根据实际需求添加
	Public.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})

	userApi := new(v1.PublicApi)
	// 用户登录与注册
	Public.POST("/user/login", userApi.Login)
	Public.POST("/user/register", userApi.Register)
	// 可继续补充如验证码、邮箱等公开接口
	// Public.POST("/user/send_captcha", userApi.SendCaptcha)
}
