// 使用 curl 测试 AI 对话接口：
// curl -X POST -H "Content-Type: application/json" -d '{"message":"你好"}' http://localhost:8080/api/public/ai/chat

package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	contextFunc "goAccounting/internal/api/util"
	aiModel "goAccounting/internal/model/ai"
)

type ChatRequest struct {
	Message   string `json:"message" binding:"required"`
	SessionId string `json:"sessionId,omitempty"`
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

	// 获取当前用户ID
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	chatService := &aiService.ChatService{}
	var responseText string
	var err error
	if req.SessionId != "" {
		responseText, err = chatService.ContinueSession(req.Message, userId, req.SessionId, ctx.Request.Context())
	} else {
		responseText, err = chatService.GetChatResponseWithUser(req.Message, userId, ctx.Request.Context())
	}
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

// 聊天历史记录请求
// GET /user/ai/chat/history?offset=0&limit=20
func GinChatHistoryHandler(ctx *gin.Context) {
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	offset := 0
	limit := 20
	if v := ctx.Query("offset"); v != "" {
		offsetInt, err := strconv.Atoi(v)
		if err == nil {
			offset = offsetInt
		}
	}
	if v := ctx.Query("limit"); v != "" {
		limitInt, err := strconv.Atoi(v)
		if err == nil {
			limit = limitInt
		}
	}
	chatService := &aiService.ChatService{}
	history, err := chatService.GetChatHistory(userId, offset, limit, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "获取历史失败: " + err.Error()})
		return
	}
	// 分组：每个sessionId只取最早一条（即每个会话的第一个问题）
	sessionMap := make(map[string]aiModel.ChatRecord)
	for i := len(history) - 1; i >= 0; i-- { // 倒序，保证最早的在后面覆盖
		record := history[i]
		sessionMap[record.SessionId] = record
	}
	var sessions []gin.H
	for _, record := range sessionMap {
		sessions = append(sessions, gin.H{
			"sessionId":     record.SessionId,
			"firstQuestion": record.Input,
			"createdAt":     record.CreatedAt,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

// 获取指定会话的全部聊天记录
// GET /user/ai/chat/session?sessionId=xxx
func GinChatSessionDetailHandler(ctx *gin.Context) {
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	sessionId := ctx.Query("sessionId")
	if sessionId == "" {
		ctx.JSON(400, gin.H{"success": false, "error": "缺少sessionId参数"})
		return
	}
	chatService := &aiService.ChatService{}
	records, err := chatService.GetSessionHistory(sessionId, 100, ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"success": false, "error": "获取会话详情失败: " + err.Error()})
		return
	}
	// 只返回属于当前用户的消息
	var filtered []aiModel.ChatRecord
	for _, r := range records {
		if r.UserId == userId || r.UserId == 0 {
			filtered = append(filtered, r)
		}
	}
	ctx.JSON(200, gin.H{"success": true, "data": filtered})
}
