package aiService

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// 建议实际项目用配置文件或环境变量
	blueLMAppID  = "2025361912"
	blueLMAppKey = "GgOxhqASWromOzvZ"
	blueLMApiUrl = "http://api-ai.vivo.com.cn/vivogpt/completions"
)

type blueLMRequest struct {
	Prompt    string `json:"prompt"`
	Model     string `json:"model"`
	SessionId string `json:"sessionId"`
}

type blueLMResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
}

type ChatService struct{}

func (s *ChatService) GetChatResponse(userInput string, ctx context.Context) (string, error) {
	sessionId := uuid.New().String()
	requestId := uuid.New().String()

	reqBody := blueLMRequest{
		Prompt:    userInput,
		Model:     "vivo-BlueLM-TB-Pro",
		SessionId: sessionId,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// 构建完整URL，包含requestId参数
	fullUrl := blueLMApiUrl + "?requestId=" + requestId

	httpReq, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// TODO: 需要添加认证头，参考 Java 代码的 VivoAuth.generateAuthHeaders

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var apiResp blueLMResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return "", err
	}

	if apiResp.Code != 0 {
		return "", errors.New(apiResp.Msg)
	}

	return apiResp.Data.Content, nil
}
