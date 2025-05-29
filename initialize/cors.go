package initialize

import (
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitCORS 初始化 CORS 配置
func InitCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",    // HBuilderX 默认端口
			"http://127.0.0.1:8080",   
			"http://localhost:3000",    // 其他可能的前端端口
			"http://127.0.0.1:3000",
			"http://localhost:8081",    // uni-app 可能的端口
			"http://127.0.0.1:8081",
			"http://localhost:5173",    // Vite 默认端口
			"http://127.0.0.1:5173",
			"http://localhost:8082",    // 其他开发端口
			"http://127.0.0.1:8082",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH",
		},
		AllowHeaders: []string{
			"Origin", 
			"Content-Type", 
			"Authorization", 
			"X-Requested-With",
			"Accept",
			"Accept-Language",
			"Content-Language",
			"Access-Control-Allow-Origin",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
