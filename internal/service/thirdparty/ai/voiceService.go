package aiService

import (
	"context"
)

type VoiceService struct{}

func (vs *VoiceService) ConvertSpeechToText(audioData []byte, ctx context.Context) (string, error) {
	// 调用第三方语音识别 API
	// 示例：伪代码
	text := "识别的文本内容"
	return text, nil
}
