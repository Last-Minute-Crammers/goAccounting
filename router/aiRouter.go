package router

import (
	"fmt"
	"goAccounting/internal/api/aiAPI"

	"github.com/gin-gonic/gin"
)

// HandlerFunc
func RegisterAIRoutes() {
	// 语音处理接口，暂未接入大模型
	PublicRouterGroup.POST("/ai/voice", func(ctx *gin.Context) {
		aiAPI.GinVoiceInputHandler(ctx)
	})
	// OCR 处理接口,暂未完全接入
	PublicRouterGroup.POST("/ai/ocr", func(ctx *gin.Context) {
		aiAPI.GinOCRInputHandler(ctx)
	})
	// AI 聊天接口，迁移到私人路由，需要认证
	Private.POST("/user/ai/chat", aiAPI.GinChatHandler)
	// AI 聊天历史记录接口
	Private.GET("/user/ai/chat/history", aiAPI.GinChatHistoryHandler)
	// AI 聊天会话详情接口
	Private.GET("/user/ai/chat/session", aiAPI.GinChatSessionDetailHandler)
	// AI 财务报告生成接口
	Private.POST("user/report/ai", aiAPI.GinAIReportHandler)
	// AI 财务报告历史接口
	Private.GET("/user/ai/history", aiAPI.GetHistoryReportHandler)
	fmt.Println("AI路由注册完成: /user/ai/history")
	
	// 添加一个测试路由，用于验证路由是否正常工作
	Private.GET("/user/ai/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "AI路由测试成功", "status": "ok"})
	})
}

