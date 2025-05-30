package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("router.engine: start init engine")
	// 这里只做一些基础配置，实际的 Engine 在 router.go 中定义
}
