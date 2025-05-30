package v1

import (
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	categoryService "goAccounting/internal/service/category"
	"goAccounting/global/constant"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryApi struct{}

var categoryServiceInstance = categoryService.NewCategoryService()

// GetList 获取分类列表
func (c *CategoryApi) GetList(ctx *gin.Context) {
	incomeExpenseStr := ctx.Query("income_expense")
	var ie *constant.IncomeExpense
	
	if incomeExpenseStr != "" {
		ieValue := constant.IncomeExpense(incomeExpenseStr)
		ie = &ieValue
	}
	
	categories, err := categoryServiceInstance.List(ie, ctx)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	var responseData []response.CategoryOne
	for _, cat := range categories {
		var catResponse response.CategoryOne
		catResponse.SetData(cat)
		responseData = append(responseData, catResponse)
	}
	
	response.OkWithData(responseData, ctx)
}

// GetOne 获取单个分类
func (c *CategoryApi) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	category, err := categoryServiceInstance.GetCategoryById(uint(id), ctx)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	var responseData response.CategoryOne
	responseData.SetData(category)
	
	response.OkWithData(responseData, ctx)
}

// CreateOne 创建分类
func (c *CategoryApi) CreateOne(ctx *gin.Context) {
	var requestData request.CategoryCreateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	category, err := categoryServiceInstance.Create(requestData.Name, requestData.IncomeExpense)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	var responseData response.CategoryOne
	responseData.SetData(category)
	
	response.OkWithData(responseData, ctx)
}

// UpdateOne 更新分类
func (c *CategoryApi) UpdateOne(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	var requestData request.CategoryUpdateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	err = categoryServiceInstance.Update(uint(id), requestData.Name, requestData.IncomeExpense, ctx)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	// 返回更新后的分类
	category, err := categoryServiceInstance.GetCategoryById(uint(id), ctx)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	var responseData response.CategoryOne
	responseData.SetData(category)
	
	response.OkWithData(responseData, ctx)
}

// DeleteOne 删除分类
func (c *CategoryApi) DeleteOne(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	err = categoryServiceInstance.Delete(uint(id), ctx)
	if err != nil {
		response.FailToInternalError(ctx, err)
		return
	}
	
	response.OkWithMessage("分类删除成功", ctx)
}
