package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine = Engine
var (
	Public  *gin.RouterGroup
	Private *gin.RouterGroup

	NoTourist *gin.RouterGroup
)

func init() {
	fmt.Println("starting init PUBLIC router")
	//Public = engine.Group(global.Config.System.RouterPrefix + "/public")
	//fmt.Println("starting init PRIVATE router")
	//Private = engine.Group(global.Config.System.RouterPrefix, middleware.JWTAuth())
	fmt.Println("init basic router success")
}
