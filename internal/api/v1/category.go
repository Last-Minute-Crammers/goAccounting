package v1

import (
	"goAccounting/global/constant"
	"goAccounting/internal/api/response"

	"github.com/gin-gonic/gin"
)

type CategoryApi struct{}

// Use the existing categoryService declared in another file

// CreateCategory
func (c *CategoryApi) CreateCategory(ctx *gin.Context) {
	var req struct {
		Name          string                 `json:"name" binding:"required"`
		IncomeExpense constant.IncomeExpense `json:"income_expense" binding:"required"`
		Icon          string                 `json:"icon" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	category, err := categoryService.Create(req.Name, req.IncomeExpense, req.Icon, ctx)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(category, ctx)
}

// UpdateCategory
func (c *CategoryApi) UpdateCategory(ctx *gin.Context) {
	var req struct {
		Name          string                 `json:"name"`
		IncomeExpense constant.IncomeExpense `json:"income_expense"`
	}
	id, ok := contextFunc.GetParamId(ctx)
	if !ok {
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	err := categoryService.Update(id, req.Name, req.IncomeExpense, ctx)
	if responseError(err, ctx) {
		return
	}
	response.OkWithMessage("更新成功", ctx)
}

// DeleteCategory
func (c *CategoryApi) DeleteCategory(ctx *gin.Context) {
	id, ok := contextFunc.GetParamId(ctx)
	if !ok {
		return
	}
	err := categoryService.Delete(id, ctx)
	if responseError(err, ctx) {
		return
	}
	response.OkWithMessage("删除成功", ctx)
}

// ListCategory
func (c *CategoryApi) ListCategory(ctx *gin.Context) {
	var req struct {
		IncomeExpense *constant.IncomeExpense `json:"income_expense" form:"income_expense"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	list, err := categoryService.List(req.IncomeExpense, ctx)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(list, ctx)
}
