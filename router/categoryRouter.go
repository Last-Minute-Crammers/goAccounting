package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterCategoryRoutes() {
	categoryApi := new(v1.CategoryApi)
	Private.POST("/category", categoryApi.CreateCategory)
	Private.PUT("/category/:id", categoryApi.UpdateCategory)
	Private.DELETE("/category/:id", categoryApi.DeleteCategory)
	Private.GET("/category/list", categoryApi.ListCategory)
}
