package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.FormValue("input")
	chatService := aiService.ChatService{}
	response, err := chatService.GetChatResponse(userInput, r.Context())
	if err != nil {
		http.Error(w, "对话失败", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("智能体回复: %s", response)))
}

// Gin 适配器
func GinChatHandler(ctx *gin.Context) {
	var req struct {
		Input string `json:"input" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	chatService := aiService.ChatService{}
	response, err := chatService.GetChatResponse(req.Input, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "对话失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": response})
}
