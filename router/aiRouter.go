package router

import (
	"goAccounting/internal/api/aiAPI"

	"github.com/gin-gonic/gin"
)

// 适配 Gin 的 HandlerFunc
func RegisterAIRoutes() {
	Private.POST("/ai/voice", func(ctx *gin.Context) {
		aiAPI.GinVoiceInputHandler(ctx)
	})
	Private.POST("/ai/ocr", func(ctx *gin.Context) {
		aiAPI.GinOCRInputHandler(ctx)
	})
	Private.POST("/ai/chat", func(ctx *gin.Context) {
		aiAPI.GinChatHandler(ctx)
	})
}
