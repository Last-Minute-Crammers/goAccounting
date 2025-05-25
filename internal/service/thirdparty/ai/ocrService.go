package aiService

import (
	"context"
)

type OCRService struct{}

func (os *OCRService) ExtractTextFromImage(imageData []byte, ctx context.Context) (string, error) {
	// 调用第三方 OCR API
	text := "提取的文字内容"
	return text, nil
}
