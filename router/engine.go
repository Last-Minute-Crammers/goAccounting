package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {
	fmt.Println("router.engine: start init engine")
	Engine = gin.Default()
	if Engine == nil {
		fmt.Println("router.engine: engine is NIL")
	}
}
