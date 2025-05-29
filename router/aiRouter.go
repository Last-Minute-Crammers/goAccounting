package router

import (
	"goAccounting/internal/api/aiAPI"
	"net/http"
)

func RegisterAIRoutes() {
	http.HandleFunc("/api/voice", aiAPI.VoiceInputHandler)
	http.HandleFunc("/api/ocr", aiAPI.OCRInputHandler)
	http.HandleFunc("/api/chat", aiAPI.ChatHandler)
}
