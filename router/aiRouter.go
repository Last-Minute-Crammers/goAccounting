package router

import (
    "net/http"
    "goAccounting/internal/api"
)

func RegisterAIRoutes() {
    http.HandleFunc("/api/voice", api.VoiceInputHandler)
    http.HandleFunc("/api/ocr", api.OCRInputHandler)
    http.HandleFunc("/api/chat", api.ChatHandler)
}