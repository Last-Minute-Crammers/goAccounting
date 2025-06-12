package router

import (
	"goAccounting/internal/api/aiAPI"

	"github.com/gin-gonic/gin"
)

// HandlerFunc
func RegisterAIRoutes() {
	// 语音处理接口，暂未接入大模型
	Public.POST("/ai/voice", func(ctx *gin.Context) {
		aiAPI.GinVoiceInputHandler(ctx)
	})
	// OCR 处理接口,暂未完全接入
	Public.POST("/ai/ocr", func(ctx *gin.Context) {
		aiAPI.GinOCRInputHandler(ctx)
	})
	// AI 聊天接口，接入大模型
	Public.POST("/ai/chat", aiAPI.GinChatHandler)
	
	// AI理财宠物相关接口
	Private.GET("/ai/pet/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Pet status retrieved",
			"data": gin.H{
				"name": "理财小汪",
				"level": 3,
				"emotion": "很开心",
				"exp": 750,
				"next_level_exp": 1000,
				"appearance": "evolved_cat",
			},
		})
	})
	
	Private.POST("/ai/pet/interact", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Pet interaction completed",
			"response": "主人今天的消费很棒哦！继续保持这个节奏～",
			"emotion_change": "happy",
			"exp_gained": 10,
		})
	})
	
	Private.GET("/ai/pet/feedback", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Daily pet feedback",
			"feedback": "今天花了150元，比昨天少了20元，小财神很开心！",
			"recommendation": "建议明天的预算控制在100元以内",
		})
	})
	
	// AI智能报告生成
	Private.POST("/ai/report/generate", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "AI report generated",
			"data": gin.H{
				"period": "monthly",
				"total_income": 8500,
				"total_expense": 6200,
				"savings_rate": 27.1,
				"spending_analysis": gin.H{
					"top_category": "餐饮",
					"trend": "increasing",
					"recommendation": "建议减少外卖频次，每周可节省约200元",
				},
				"financial_health_score": 78,
				"user_type": "稳健理财型",
			},
		})
	})
	
	Private.GET("/ai/report/insights", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Financial insights",
			"insights": []gin.H{
				{"type": "warning", "message": "本周外卖消费增长60%"},
				{"type": "tip", "message": "建议使用52周存钱法，年底可存款13780元"},
				{"type": "achievement", "message": "连续7天控制预算内消费，获得'自律达人'徽章"},
			},
		})
	})
	
	// AI理财建议
	Private.GET("/ai/advice/personalized", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Personalized financial advice",
			"advice": []gin.H{
				{"category": "支出优化", "suggestion": "减少娱乐支出15%，可增加月储蓄500元"},
				{"category": "投资建议", "suggestion": "当前储蓄率良好，可考虑定投基金"},
				{"category": "消费习惯", "suggestion": "工作日午餐平均45元，建议控制在35元"},
			},
		})
	})
}

