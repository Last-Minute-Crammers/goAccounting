package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterTransactionRoutes() {
	if NoTourist == nil {
		return
	}
	
	transactionApi := v1.TransactionApi{}
	
	transactionGroup := NoTourist.Group("/transaction")
	{
		// 基础CRUD
		transactionGroup.GET("/list", transactionApi.GetList)
		transactionGroup.POST("", transactionApi.CreateOne)
		transactionGroup.GET("/:id", transactionApi.GetOne)
		
		// 统计相关
		statisticGroup := transactionGroup.Group("/statistic")
		{
			statisticGroup.GET("/month", transactionApi.GetMonthStatistic)
		}
	}
}
