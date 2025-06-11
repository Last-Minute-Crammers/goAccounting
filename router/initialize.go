package router

import (
	"fmt"
	"goAccounting/global"
	"goAccounting/router/middleware"

	"github.com/gin-gonic/gin"
)

var (
	PublicRouterGroup *gin.RouterGroup
	Private           *gin.RouterGroup
	NoTourist         *gin.RouterGroup
)

// 添加占位符路由注册函数
func RegisterTransactionRoutesPlaceholder() {
	// TODO: 实现交易相关路由
	Private.GET("/transaction/list", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Transaction routes not implemented yet"})
	})
}

func RegisterCategoryRoutes() {
	RegisterCategoryRoutesV1()
}

// Public/Private 路由组初始化后可用于注册路由
func init() {
	fmt.Println("starting init PUBLIC router")
	if Engine == nil {
		fmt.Println("router: engine is NIL")
		return
	}

	// 初始化路由组
	PublicRouterGroup = Engine.Group(global.Config.System.RouterPrefix + "/public")
	fmt.Println("starting init PRIVATE router")
	Private = Engine.Group(global.Config.System.RouterPrefix, middleware.JWTAuth())
	fmt.Println("init basic router success")

	// 注册所有路由 - 先注册Public路由，再注册AI路由
	RegisterPublicRoutes()
	RegisterAIRoutes() // AI routes now use Public group
	RegisterUserRoutes()
		RegisterTransactionRoutesPlaceholder()
	RegisterCategoryRoutes()
	RegisterTransactionRoutes()

	fmt.Println("All routes registered successfully")
}
