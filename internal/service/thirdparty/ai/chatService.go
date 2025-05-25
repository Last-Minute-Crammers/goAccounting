package aiService

import (
	"context"
)

type ChatService struct{}

func (cs *ChatService) GetChatResponse(userInput string, ctx context.Context) (string, error) {
	// 调用第三方对话 API
	response := "智能体的回复内容"
	return response, nil
}
