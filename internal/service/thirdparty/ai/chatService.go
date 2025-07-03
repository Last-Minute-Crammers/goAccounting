package aiService

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"goAccounting/global/db"
	aiModel "goAccounting/internal/model/ai"
	"github.com/google/uuid"
)

const (
	// 根据官方文档的配置
	blueLMAppID  = "2025361912"
	blueLMAppKey = "GgOxhqASWromOzvZ"
	blueLMApiUrl = "https://api-ai.vivo.com.cn/vivogpt/completions"
)

// 按照官方文档的请求体结构
type blueLMRequest struct {
	Prompt       string                 `json:"prompt"`
	Model        string                 `json:"model"`
	SessionId    string                 `json:"sessionId"`
	SystemPrompt string                 `json:"systemPrompt,omitempty"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
}

// 按照官方文档的响应结构
type blueLMResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SessionId string `json:"sessionId"`
		RequestId string `json:"requestId"`
		Content   string `json:"content"`
		Usage     struct {
			PromptTokens     int `json:"promptTokens"`
			CompletionTokens int `json:"completionTokens"`
			TotalTokens      int `json:"totalTokens"`
		} `json:"usage"`
		FinishReason string `json:"finishReason"`
	} `json:"data"`
}

type ChatService struct{}

// GetChatResponse 获取AI聊天响应并保存记录
func (s *ChatService) GetChatResponse(userInput string, ctx context.Context) (string, error) {
	return s.GetChatResponseWithUser(userInput, 0, ctx) // 默认用户ID为0
}

// GetChatResponseWithUser 获取AI聊天响应并保存记录（带用户ID）
func (s *ChatService) GetChatResponseWithUser(userInput string, userId uint, ctx context.Context) (string, error) {
	return s.GetChatResponseWithSession(userInput, userId, "", ctx) // 创建新会话
}

// GetChatResponseWithSession 获取AI聊天响应并保存记录（指定会话ID）
func (s *ChatService) GetChatResponseWithSession(userInput string, userId uint, sessionId string, ctx context.Context) (string, error) {
	// 如果没有提供会话ID，创建新会话
	if sessionId == "" {
		sessionId = uuid.New().String()
	}
	requestId := uuid.New().String()

	log.Printf("=== 开始AI对话请求 ===")
	log.Printf("RequestID: %s", requestId)
	log.Printf("SessionID: %s", sessionId)
	log.Printf("UserID: %d", userId)
	log.Printf("用户输入: %s", userInput)

	// 获取用户理财统计数据（示例，实际可根据业务调整）
	financialData := ""
	// TODO: 调用statisticService获取用户理财统计数据，拼接为字符串
	// 例如: financialData = GetUserFinancialSummary(userId, ctx)

	// 获取整理后的对话历史
	historySummary := ""
	records, err := s.GetSessionHistory(sessionId, 1, ctx) // 获取上一轮对话
	if err == nil && len(records) > 0 {
		historySummary = records[len(records)-1].HistorySummary
	}

	// 动态生成prompt
	systemPrompt := "你的中文名字叫理财小汪，你是智能理财宠物，对用户的称呼是主人。你有着丰富的理财知识，活泼可爱，认真可靠，你需要协助用户进行个人理财规划。当回复问题时需要回复你的名字时，中文名必须回复理财小汪，此外回复和你的名字相关的问题时，也需要给出和你的名字对应的合理回复。"
	fullPrompt := fmt.Sprintf("%s\n\n用户的理财数据:\n%s\n\n对话历史:\n%s\n\n用户提问:\n%s", systemPrompt, financialData, historySummary, userInput)

	// 创建聊天记录，在服务层完全控制ID
	chatRecord := &aiModel.ChatRecord{
		SessionId:      sessionId,
		RequestId:      requestId,
		UserId:         userId,
		Input:          userInput,
		Response:       "", // 先创建记录，响应后更新
		HistorySummary: historySummary, // 保存整理后的对话历史
	}

	// 获取AI响应
	response, err := s.callBlueLMAPI(fullPrompt, sessionId, requestId, ctx)
	if err != nil {
		log.Printf("AI API调用失败: %v", err)
		// 即使API调用失败，也保存记录，标记错误
		chatRecord.Response = fmt.Sprintf("API调用失败: %v", err)
		s.saveChatRecord(chatRecord, ctx)
		return "", err
	}

	// 更新响应内容
	chatRecord.Response = response

	// TODO: 这里可以调用大模型对本轮对话进行总结，更新HistorySummary
	// chatRecord.HistorySummary = UpdateHistorySummary(historySummary, userInput, response)

	// 保存聊天记录到数据库
	if err := s.saveChatRecord(chatRecord, ctx); err != nil {
		log.Printf("保存聊天记录失败: %v", err)
		// 不影响主流程，继续返回响应
	}

	log.Printf("=== AI对话成功完成 ===")
	return response, nil
}

// callBlueLMAPI 调用蓝心大模型API
func (s *ChatService) callBlueLMAPI(userInput, sessionId, requestId string, ctx context.Context) (string, error) {
	// 定义系统提示词
	systemPrompt := "你的中文名字叫理财小汪，你是智能理财宠物，对用户的称呼是主人。你有着丰富的理财知识，活泼可爱，认真可靠，你需要协助用户进行个人理财规划。当回复问题时需要回复你的名字时，中文名必须回复理财小汪，此外回复和你的名字相关的问题时，也需要给出和你的名字对应的合理回复。"

	// 按照Python示例构建请求体，添加系统提示词
	reqBody := blueLMRequest{
		Prompt:       userInput,
		Model:        "vivo-BlueLM-TB-Pro", // 使用与Python示例一致的模型
		SessionId:    sessionId,
		SystemPrompt: systemPrompt,
		Extra: map[string]interface{}{
			"temperature":  0.9,
			"systemPrompt": systemPrompt, // 在extra中也添加一份，确保兼容性
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("请求体序列化失败: %v", err)
	}

	log.Printf("请求体JSON: %s", string(bodyBytes))
	log.Printf("系统提示词: %s", systemPrompt)

	// 准备查询参数
	queryParams := map[string]string{
		"requestId": requestId,
	}

	// 构建URL
	fullUrl := blueLMApiUrl
	log.Printf("请求URL: %s", fullUrl)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fullUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 添加查询参数
	q := httpReq.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 按照Python示例生成认证头
	authHeaders := GenerateAuthHeaders("POST", "/vivogpt/completions", queryParams, blueLMAppID, blueLMAppKey)
	log.Printf("认证头信息:")
	for key, value := range authHeaders {
		log.Printf("  %s: %s", key, value)
		httpReq.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	log.Printf("发送HTTP请求到蓝心大模型API...")
	startTime := time.Now()
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("HTTP请求失败: %v", err)
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)
	log.Printf("HTTP请求完成，耗时: %v", duration)
	log.Printf("HTTP响应状态码: %d", resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	log.Printf("原始响应内容 (长度: %d): %s", len(respBytes), string(respBytes))

	// 如果状态码不是200，直接返回错误
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP状态码错误: %d, 响应: %s", resp.StatusCode, string(respBytes))
	}

	var apiResp blueLMResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		log.Printf("JSON解析失败: %v", err)
		return "", fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(respBytes))
	}

	log.Printf("=== API响应解析结果 ===")
	log.Printf("Code: %d", apiResp.Code)
	log.Printf("Msg: %s", apiResp.Msg)
	log.Printf("Content: %s", apiResp.Data.Content)

	if apiResp.Code != 0 {
		log.Printf("API返回错误码: %d, 错误信息: %s", apiResp.Code, apiResp.Msg)
		return "", fmt.Errorf("API调用失败: %s (code: %d)", apiResp.Msg, apiResp.Code)
	}

	// 检查内容是否为空
	if apiResp.Data.Content == "" {
		log.Printf("警告: API返回的Content字段为空")
		return "", nil
	}

	return apiResp.Data.Content, nil
}

// saveChatRecord 保存聊天记录到数据库
func (s *ChatService) saveChatRecord(record *aiModel.ChatRecord, ctx context.Context) error {
	dao := aiModel.NewChatDAO(db.GetDb(ctx))
	if err := dao.Create(record); err != nil {
		log.Printf("保存聊天记录失败: %v", err)
		return err
	}
	log.Printf("聊天记录保存成功 - RequestID: %s", record.RequestId)
	return nil
}

// GetChatHistory 获取用户聊天历史
func (s *ChatService) GetChatHistory(userId uint, offset, limit int, ctx context.Context) ([]aiModel.ChatRecord, error) {
	dao := aiModel.NewChatDAO(db.GetDb(ctx))
	return dao.GetByUserId(userId, offset, limit)
}

// GetSessionHistory 获取指定会话的聊天历史
func (s *ChatService) GetSessionHistory(sessionId string, limit int, ctx context.Context) ([]aiModel.ChatRecord, error) {
	dao := aiModel.NewChatDAO(db.GetDb(ctx))
	return dao.GetBySessionId(sessionId, limit)
}

// CreateNewSession 创建新的会话ID
func (s *ChatService) CreateNewSession() string {
	return uuid.New().String()
}

// ContinueSession 继续现有会话
func (s *ChatService) ContinueSession(userInput string, userId uint, sessionId string, ctx context.Context) (string, error) {
	return s.GetChatResponseWithSession(userInput, userId, sessionId, ctx)
}
