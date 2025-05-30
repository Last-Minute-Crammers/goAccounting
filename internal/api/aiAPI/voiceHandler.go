package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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

// Gin 适配器
func GinVoiceInputHandler(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("audio")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "音频文件获取失败"})
		return
	}
	defer file.Close()
	audioData, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "音频读取失败"})
		return
	}
	voiceService := aiService.VoiceService{}
	text, err := voiceService.ConvertSpeechToText(audioData, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "语音识别失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": text})
}
