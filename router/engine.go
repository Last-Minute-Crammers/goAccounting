package router

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {
	fmt.Println("router.engine: start init engine")
	Engine = gin.Default()
	// 添加自定义 CORS 中间件，允许所有来源和常用方法
	Engine.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false, 
		MaxAge:           12 * 60 * 60, 
	}))


	Engine.Use(func(c *gin.Context) {
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
			c.Header("Access-Control-Max-Age", "43200") // 12小时
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 添加测试路由，不需要认证
	Engine.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "测试路由正常", "status": "ok"})
	})

	if Engine == nil {
		fmt.Println("router.engine: engine is NIL")
	} else {
		fmt.Println("router.engine: engine initialized successfully")
	}
}
