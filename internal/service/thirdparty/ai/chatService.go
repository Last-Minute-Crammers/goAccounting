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

func (s *ChatService) GetChatResponse(userInput string, ctx context.Context) (string, error) {
	sessionId := uuid.New().String()
	requestId := uuid.New().String()

	log.Printf("=== 开始AI对话请求 ===")
	log.Printf("RequestID: %s", requestId)
	log.Printf("SessionID: %s", sessionId)
	log.Printf("用户输入: %s", userInput)

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

	// 打印完整请求信息
	log.Printf("完整请求头:")
	for key, values := range httpReq.Header {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

	log.Printf("最终请求URL: %s", httpReq.URL.String())

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
	log.Printf("HTTP响应状态: %s", resp.Status)

	// 打印响应头
	log.Printf("响应头:")
	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

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
	log.Printf("SessionId: %s", apiResp.Data.SessionId)
	log.Printf("RequestId: %s", apiResp.Data.RequestId)
	log.Printf("Content: %s", apiResp.Data.Content)
	log.Printf("Content长度: %d", len(apiResp.Data.Content))
	log.Printf("FinishReason: %s", apiResp.Data.FinishReason)
	log.Printf("Usage: %+v", apiResp.Data.Usage)

	if apiResp.Code != 0 {
		log.Printf("API返回错误码: %d, 错误信息: %s", apiResp.Code, apiResp.Msg)
		return "", fmt.Errorf("API调用失败: %s (code: %d)", apiResp.Msg, apiResp.Code)
	}

	// 检查内容是否为空
	if apiResp.Data.Content == "" {
		log.Printf("警告: API返回的Content字段为空")
		return "", nil
	}

	log.Printf("=== AI对话成功 ===")
	log.Printf("返回内容长度: %d 字符", len(apiResp.Data.Content))
	return apiResp.Data.Content, nil
}
