package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterPublicRoutes() {
	if Public == nil {
		return
	}
	
	publicApi := v1.PublicApi{}
	
	userGroup := Public.Group("/user")
	{
		userGroup.POST("/login", publicApi.Login)
		userGroup.POST("/register", publicApi.Register)
	}
	
	// 健康检查
	Public.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
