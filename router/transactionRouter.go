package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterTransactionRoutes() {
	transactionApi := new(v1.TransactionApi)
	Private.GET("/user/transaction/:id", transactionApi.GetOne)
	Private.POST("/user/transaction/list", transactionApi.GetList) // 改为POST支持JSON body
	Private.POST("/user/transaction", transactionApi.CreateOne)
	Private.POST("/user/transaction/statistic/day", transactionApi.GetDayStatistic)     // 改为POST支持JSON body
	Private.POST("/user/transaction/statistic/week", transactionApi.GetWeekStatistic)   // 改为POST支持JSON body
	Private.POST("/user/transaction/statistic/month", transactionApi.GetMonthStatistic) // 改为POST支持JSON body
	Private.POST("/user/transaction/statistic/year", transactionApi.GetYearStatistic)   // 改为POST支持JSON body
	Private.GET("/user/transaction/statistic/total", transactionApi.GetTotalStatistic)

	Private.GET("/user/transaction/statistic/category_rank", transactionApi.GetCategoryAmountRank)
	Private.GET("/user/transaction/statistic/amount_rank", transactionApi.GetAmountRank)
}
