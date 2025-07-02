package aiService

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ocrAppID  = "2025361912"  // Replace with actual App ID
	ocrAppKey = "GgOxhqASWromOzvZ" // Replace with actual App Key
	ocrApiUrl = "https://api-ai.vivo.com.cn/ocr/extractText"
)

type ocrRequest struct {
	ImageBase64 string `json:"imageBase64"`
}

type ocrResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

type OCRService struct{}

func (os *OCRService) ExtractTextFromImage(imageData []byte, ctx context.Context) (string, error) {
	// Convert image data to Base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// Prepare request body
	reqBody := ocrRequest{
		ImageBase64: imageBase64,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("请求体序列化失败: %v", err)
	}

	// Prepare HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", ocrApiUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Generate authentication headers
	authHeaders := GenerateAuthHeaders("POST", "/ocr/extractText", nil, ocrAppID, ocrAppKey)
	for key, value := range authHeaders {
		httpReq.Header.Set(key, value)
	}

	// Send HTTP request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP状态码错误: %d, 响应: %s", resp.StatusCode, string(respBytes))
	}

	var apiResp ocrResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return "", fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(respBytes))
	}

	if apiResp.Code != 0 {
		return "", fmt.Errorf("OCR API调用失败: %s (code: %d)", apiResp.Msg, apiResp.Code)
	}

	return apiResp.Data.Text, nil
}
