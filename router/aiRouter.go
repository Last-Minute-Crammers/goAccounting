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
}

