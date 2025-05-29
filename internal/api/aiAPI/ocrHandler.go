package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"net/http"
)

func OCRInputHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取图片数据
	imageData := []byte{} // 示例：从请求中解析图片数据
	ocrService := aiService.OCRService{}
	text, err := ocrService.ExtractTextFromImage(imageData, r.Context())
	if err != nil {
		http.Error(w, "OCR 识别失败", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("提取结果: %s", text)))
}
