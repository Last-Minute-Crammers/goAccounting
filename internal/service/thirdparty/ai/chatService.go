package aiService

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// 建议实际项目用配置文件或环境变量
	blueLMAppID  = "2025361912"
	blueLMAppKey = "GgOxhqASWromOzvZ"
	blueLMApiUrl = "https://open.vivo.com.cn/openapi/llm/v1/chat/completions"
)

type blueLMRequest struct {
	AppID      string `json:"appId"`
	AppKey     string `json:"appKey"`
	Messages   []blueLMMessage `json:"messages"`
	Model      string `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type blueLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type blueLMResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Choices []struct {
			Message blueLMMessage `json:"message"`
		} `json:"choices"`
	} `json:"data"`
}

type ChatService struct{}

func (s *ChatService) GetChatResponse(userInput string, ctx context.Context) (string, error) {
	reqBody := blueLMRequest{
		AppID:  blueLMAppID,
		AppKey: blueLMAppKey,
		Messages: []blueLMMessage{
			{Role: "user", Content: userInput},
		},
		Model: "blueLM-32B-Chat", // 可根据文档调整
		Temperature: 0.7,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	httpReq, err := http.NewRequest("POST", blueLMApiUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
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
	if len(apiResp.Data.Choices) == 0 {
		return "", errors.New("no response from 蓝心大模型")
	}
	return apiResp.Data.Choices[0].Message.Content, nil
}
