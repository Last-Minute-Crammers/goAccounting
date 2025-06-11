package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterTransactionRoutes() {
	transactionApi := new(v1.TransactionApi)
	Private.GET("/user/transaction/:id", transactionApi.GetOne)
	Private.GET("/user/transaction/list", transactionApi.GetList)
	Private.POST("/user/transaction", transactionApi.CreateOne)
	Private.GET("/user/transaction/statistic/month", transactionApi.GetMonthStatistic)
	Private.GET("/user/transaction/statistic/total", transactionApi.GetTotalStatistic)
}
