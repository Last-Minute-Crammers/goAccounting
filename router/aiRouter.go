package router

import (
	"goAccounting/internal/api/aiAPI"

	"github.com/gin-gonic/gin"
)

// 适配 Gin 的 HandlerFunc - 改为公开路由用于测试
func RegisterAIRoutes() {
	// 使用 Public 路由组，不需要 JWT 认证
	Public.POST("/ai/voice", func(ctx *gin.Context) {
		aiAPI.GinVoiceInputHandler(ctx)
	})
	Public.POST("/ai/ocr", func(ctx *gin.Context) {
		aiAPI.GinOCRInputHandler(ctx)
	})
	Public.POST("/ai/chat", func(ctx *gin.Context) {
		aiAPI.GinChatHandler(ctx)
	})
}
