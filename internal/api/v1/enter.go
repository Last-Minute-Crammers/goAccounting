package v1

import (
	"goAccounting/internal/api/response"
	apiUtil "goAccounting/internal/api/util"
	"goAccounting/internal/service"

	"github.com/gin-gonic/gin"
)

type PublicApi struct{}

var contextFunc = apiUtil.ContextFunc
var checkFunc = apiUtil.CheckFunc

var (
	commonService = service.GroupApp.CommonServiceGroup
)
var (
	userService        = service.GroupApp.UserServiceGroup
	categoryService    = service.GroupApp.CategoryServiceGroup
	transactionService = service.GroupApp.TransactionServiceGroup
)

func responseError(err error, ctx *gin.Context) bool {
	if err != nil {
		response.FailToError(ctx, err)
		return true
	}
	return false
}
