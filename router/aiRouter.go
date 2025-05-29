package router

import (
	"goAccounting/router/websocket"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func RegisterAIRoutes() {
	if NoTourist == nil {
		return
	}
	
	aiGroup := NoTourist.Group("/ai")
	{
		// 语音输入处理
		aiGroup.POST("/voice", func(c *gin.Context) {
			// TODO: 实现语音处理逻辑
			c.JSON(http.StatusOK, gin.H{
				"text": "识别到的文本内容",
				"confidence": 0.95,
			})
		})
		
		// OCR识别处理
		aiGroup.POST("/ocr", func(c *gin.Context) {
			// TODO: 实现OCR识别逻辑
			c.JSON(http.StatusOK, gin.H{
				"text": "识别到的文本内容",
				"amount": 50.00,
				"category": "餐饮",
			})
		})
		
		// AI聊天接口
		aiGroup.POST("/chat", func(c *gin.Context) {
			// TODO: 实现AI聊天逻辑
			c.JSON(http.StatusOK, gin.H{
				"reply": "您好，我是您的AI理财助手，很高兴为您服务！",
			})
		})
		
		// WebSocket聊天连接
		aiGroup.GET("/chat/ws", websocket.UpgradeToWebsocket(func(conn *websocket.Conn, ctx *gin.Context) error {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				// 回显消息或调用AI处理
				err = conn.WriteMessage(websocket.TextMessage, []byte("AI回复: "+string(message)))
				if err != nil {
					break
				}
			}
			return nil
		}))
		
		// 获取AI理财建议
		aiGroup.GET("/suggestions", func(c *gin.Context) {
			// TODO: 实现获取AI建议逻辑
			c.JSON(http.StatusOK, gin.H{
				"suggestions": []gin.H{
					{
						"type": "减少支出",
						"title": "减少外卖开支", 
						"description": "您本月外卖支出较高，建议适当减少",
						"potential_saving": 500,
					},
				},
			})
		})
	}
}
