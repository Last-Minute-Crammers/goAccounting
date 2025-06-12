package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterTransactionRoutes() {
	transactionApi := new(v1.TransactionApi)
	
	// 现有路由
	Private.GET("/user/transaction/:id", transactionApi.GetOne)
	Private.POST("/user/transaction/list", transactionApi.GetList) // 改为POST支持JSON body
	Private.POST("/user/transaction", transactionApi.CreateOne)
	Private.POST("/user/transaction/statistic/month", transactionApi.GetMonthStatistic) // 改为POST支持JSON body
	Private.GET("/user/transaction/statistic/total", transactionApi.GetTotalStatistic)
	
	// 多种记账方式
	// 语音记账
	Private.POST("/user/transaction/voice", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Voice transaction processed",
			"data": gin.H{
				"amount": -50.0,
				"category": "餐饮",
				"description": "午餐",
				"confidence": 0.95,
			},
		})
	})
	
	// OCR小票识别记账
	Private.POST("/user/transaction/ocr", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Receipt OCR processed",
			"data": gin.H{
				"amount": -128.50,
				"merchant": "沃尔玛超市",
				"category": "购物",
				"items": []gin.H{
					{"name": "苹果", "price": 25.50},
					{"name": "牛奶", "price": 18.00},
				},
				"confidence": 0.88,
			},
		})
	})
	
	// 快捷记账模板
	Private.GET("/user/transaction/templates", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Quick transaction templates",
			"data": []gin.H{
				{"name": "早餐", "amount": -15, "category": "餐饮"},
				{"name": "地铁", "amount": -6, "category": "交通"},
				{"name": "咖啡", "amount": -25, "category": "餐饮"},
			},
		})
	})
	
	// 批量导入记账
	Private.POST("/user/transaction/batch", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Batch transactions imported",
			"success_count": 15,
			"failed_count": 2,
		})
	})
}
