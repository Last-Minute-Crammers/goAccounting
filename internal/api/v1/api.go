package v1

import (
	"goAccounting/internal/api/response"
	"goAccounting/internal/api/util"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	contextFunc = util.ContextFunc
)

// PublicApi 公开接口
type PublicApi struct{}

// responseError 统一错误处理
func responseError(err error, ctx *gin.Context) bool {
	if err != nil {
		log.Printf("API Error: %v", err)
		response.FailToError(ctx, err)
		return true
	}
	return false
}
