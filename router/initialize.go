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

	NoTourist *gin.RouterGroup
)

// Public/Private 路由组初始化后可用于注册路由
func init() {
	fmt.Println("starting init PUBLIC router")
	if Engine == nil {
		fmt.Println("router: engine is NIL")
		return
	}
	PublicRouterGroup = Engine.Group(global.Config.System.RouterPrefix + "/public")
	fmt.Println("starting init PRIVATE router")
	Private = Engine.Group(global.Config.System.RouterPrefix, middleware.JWTAuth())
	fmt.Println("init basic router success")
}
