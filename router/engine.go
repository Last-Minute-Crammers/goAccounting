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
	// 添加 CORS 中间件
	Engine.Use(cors.Default())
	if Engine == nil {
		fmt.Println("router.engine: engine is NIL")
	} else {
		fmt.Println("router.engine: engine initialized successfully")
	}
}
