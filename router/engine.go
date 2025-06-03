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
		ExposeHeaders: []string{"Content-Length"},
		//AllowCredentials: true,
	}))
	if Engine == nil {
		fmt.Println("router.engine: engine is NIL")
	} else {
		fmt.Println("router.engine: engine initialized successfully")
	}
}
