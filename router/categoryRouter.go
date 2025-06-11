package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterCategoryRoutesV1() {
	categoryApi := new(v1.CategoryApi)
	Private.POST("/user/category", categoryApi.CreateCategory)
	Private.PUT("/user/category/:id", categoryApi.UpdateCategory)
	Private.DELETE("/user/category/:id", categoryApi.DeleteCategory)
	Private.GET("/user/category/list", categoryApi.ListCategory)
}
