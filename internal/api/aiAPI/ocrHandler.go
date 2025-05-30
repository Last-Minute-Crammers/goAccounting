package aiAPI

import (
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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

// Gin 适配器
func GinOCRInputHandler(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "图片文件获取失败"})
		return
	}
	defer file.Close()
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "图片读取失败"})
		return
	}
	ocrService := aiService.OCRService{}
	text, err := ocrService.ExtractTextFromImage(imageData, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OCR 识别失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": text})
}
