// 使用 curl 测试 AI 对话接口：
// curl -X POST -H "Content-Type: application/json" -d '{"message":"你好"}' http://localhost:8080/api/public/ai/chat

package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Debug   string `json:"debug,omitempty"` // 添加调试信息
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.FormValue("input")
	chatService := &aiService.ChatService{}
	responseText, err := chatService.GetChatResponse(userInput, r.Context())
	if err != nil {
		http.Error(w, "对话失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("智能体回复: %s", responseText)))
}

// Gin 适配器
func GinChatHandler(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		ctx.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "输入内容不能为空",
		})
		return
	}

	log.Printf("收到AI对话请求: %s", req.Message)

	chatService := &aiService.ChatService{}
	responseText, err := chatService.GetChatResponse(req.Message, ctx.Request.Context())
	if err != nil {
		log.Printf("AI对话失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, ChatResponse{
			Success: false,
			Error:   "对话失败: " + err.Error(),
			Debug:   fmt.Sprintf("请求内容: %s", req.Message),
		})
		return
	}

	log.Printf("AI对话成功，响应长度: %d", len(responseText))
	log.Printf("AI响应内容: %s", responseText)
	
	// 检查返回内容是否为空
	if responseText == "" {
		log.Printf("警告: AI返回了空内容")
		ctx.JSON(http.StatusOK, ChatResponse{
			Success: false,
			Error:   "AI返回了空内容",
			Debug:   "API调用成功但返回空字符串",
		})
		return
	}

	ctx.JSON(http.StatusOK, ChatResponse{
		Success: true,
		Data:    responseText,
	})
}
