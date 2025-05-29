package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"net/http"
)

func VoiceInputHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取音频数据
	audioData := []byte{} // 示例：从请求中解析音频数据
	voiceService := aiService.VoiceService{}
	text, err := voiceService.ConvertSpeechToText(audioData, r.Context())
	if err != nil {
		http.Error(w, "语音识别失败", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("识别结果: %s", text)))
}
