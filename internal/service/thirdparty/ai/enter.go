package aiService

import (
	"context"
	"fmt"
)

type ServiceGroup struct {
	VoiceService  *VoiceService
	OCRService    *OCRService
	ChatService   *ChatService
	ReportService *ReportService
	PetService    *PetService
}

func NewServiceGroup() *ServiceGroup {
	return &ServiceGroup{
		VoiceService:  &VoiceService{},
		OCRService:    &OCRService{},
		ChatService:   &ChatService{},
		ReportService: &ReportService{},
		PetService:    &PetService{},
	}
}

// AIService 统一的AI服务接口
type AIService interface {
	ProcessRequest(ctx context.Context, requestType string, data interface{}) (interface{}, error)
}

type UnifiedAIService struct {
	services *ServiceGroup
}

func NewUnifiedAIService() *UnifiedAIService {
	return &UnifiedAIService{
		services: NewServiceGroup(),
	}
}

func (uas *UnifiedAIService) ProcessRequest(ctx context.Context, requestType string, data interface{}) (interface{}, error) {
	switch requestType {
	case "voice_to_text":
		if audioData, ok := data.([]byte); ok {
			return uas.services.VoiceService.ConvertSpeechToText(audioData, ctx)
		}
	case "ocr":
		if imageData, ok := data.([]byte); ok {
			return uas.services.OCRService.ExtractTextFromImage(imageData, ctx)
		}
	case "chat":
		if userInput, ok := data.(string); ok {
			return uas.services.ChatService.GetChatResponse(userInput, ctx)
		}
	case "pet_mood":
		if goal, ok := data.(FinancialGoal); ok {
			return uas.services.PetService.UpdatePetMood(goal, ctx)
		}
	case "daily_encouragement":
		return uas.services.PetService.GetDailyEncouragement(ctx)
	}

	return nil, fmt.Errorf("unsupported request type: %s", requestType)
}
