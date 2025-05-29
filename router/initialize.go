package router

import (
	"fmt"
	"goAccounting/global"
	"goAccounting/initialize"
	"goAccounting/router/middleware"

	"github.com/gin-gonic/gin"
)

var (
	Public    *gin.RouterGroup
	Private   *gin.RouterGroup
	NoTourist *gin.RouterGroup
)

func init() {
	fmt.Println("starting init PUBLIC router")
	if Engine == nil {
		fmt.Println("router: engine is NIL")
		return
	}

	// 使用统一的CORS配置
	initialize.InitCORS(Engine)

	Public = Engine.Group(global.Config.System.RouterPrefix + "/public")
	fmt.Println("starting init PRIVATE router")
	Private = Engine.Group(global.Config.System.RouterPrefix, middleware.JWTAuth())
	fmt.Println("starting init NO TOURIST router")
	NoTourist = Engine.Group(global.Config.System.RouterPrefix, middleware.JWTAuth())

	// 注册所有路由
	RegisterPublicRoutes()
	RegisterUserRoutes()
	RegisterTransactionRoutes()
	RegisterAIRoutes()
	RegisterSocialRoutes()
	RegisterReportRoutes()

	fmt.Println("init basic router success")
}
