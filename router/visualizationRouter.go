package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterVisualizationRoutes() {
	// 数据可视化相关路由
	
	// 收支趋势图
	Private.GET("/user/chart/income-expense-trend", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Income expense trend chart data",
			"data": gin.H{
				"labels": []string{"1月", "2月", "3月", "4月", "5月", "6月"},
				"income": []float64{8500, 8200, 9000, 8800, 9200, 8600},
				"expense": []float64{6200, 5800, 6500, 6100, 6800, 6300},
				"chart_type": "line",
			},
		})
	})
	
	// 消费分类饼图
	Private.GET("/user/chart/expense-category", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Expense category pie chart data",
			"data": []gin.H{
				{"category": "餐饮", "amount": 1800, "percentage": 30},
				{"category": "交通", "amount": 900, "percentage": 15},
				{"category": "娱乐", "amount": 1200, "percentage": 20},
				{"category": "购物", "amount": 1500, "percentage": 25},
				{"category": "其他", "amount": 600, "percentage": 10},
			},
		})
	})
	
	// 预算执行情况
	Private.GET("/user/chart/budget-execution", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Budget execution chart data",
			"data": gin.H{
				"categories": []string{"餐饮", "交通", "娱乐", "购物"},
				"budgets": []float64{2000, 800, 1000, 1500},
				"actual": []float64{1800, 900, 1200, 1300},
				"chart_type": "bar",
			},
		})
	})
	
	// 消费热力图
	Private.GET("/user/chart/spending-heatmap", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Spending heatmap data",
			"data": gin.H{
				"dates": []string{"2024-01", "2024-02", "2024-03"},
				"categories": []string{"餐饮", "交通", "娱乐"},
				"values": [][]float64{
					{1800, 900, 1200},
					{1600, 850, 1100},
					{1900, 950, 1300},
				},
			},
		})
	})
	
	// 储蓄目标进度
	Private.GET("/user/chart/savings-progress", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Savings progress chart data",
			"data": gin.H{
				"target": 50000,
				"current": 32500,
				"monthly_progress": []gin.H{
					{"month": "1月", "saved": 2500},
					{"month": "2月", "saved": 2800},
					{"month": "3月", "saved": 2200},
				},
				"completion_rate": 65,
			},
		})
	})
	
	// 消费异常检测
	Private.GET("/user/analysis/anomaly", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Spending anomaly analysis",
			"anomalies": []gin.H{
				{
					"date": "2024-03-15",
					"category": "购物",
					"amount": 1200,
					"expected": 400,
					"deviation": 200,
					"description": "单日购物支出异常增高",
				},
				{
					"date": "2024-03-20",
					"category": "餐饮",
					"amount": 300,
					"expected": 80,
					"deviation": 275,
					"description": "单餐消费超出正常范围",
				},
			},
		})
	})
	
	// 理财知识推送
	Private.GET("/user/knowledge/articles", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Financial knowledge articles",
			"articles": []gin.H{
				{
					"title": "如何制定个人月度预算",
					"summary": "详细介绍月度预算制定的步骤和技巧",
					"category": "预算管理",
					"read_time": 5,
				},
				{
					"title": "延迟满足：战胜冲动消费的心理技巧",
					"summary": "学会延迟满足，减少不必要的消费支出",
					"category": "消费心理",
					"read_time": 8,
				},
			},
		})
	})
}
