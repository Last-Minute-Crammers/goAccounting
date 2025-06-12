package websocket

import (
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeToWebsocket(handler func(conn *websocket.Conn, ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade websocket connection: %v\n", err)
			ctx.JSON(500, gin.H{"error": "Failed to upgrade websocket connection"})
			return
		}
		conn.SetPingHandler(
			func(message string) error {
				err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(1))
				// if connnction is closed else if network have temperr else return err
				if err == websocket.ErrCloseSent {
					return nil
				} else if e, ok := err.(net.Error); ok && e.Temporary() {
					return nil
				}
				return err
			},
		)
		conn.SetPongHandler(nil)
		conn.SetCloseHandler(nil)
		defer conn.Close()
		err = handler(conn, ctx)
		if err != nil {
			log.Println("websocket err")
		}

	}
}

func RegisterWebsocketRoutes(router *gin.Engine) {
	Private := router.Group("/private")
	Private.GET("/ws/echo", UpgradeToWebsocket(EchoHandler))
	
	// 社交互动实时通知
	Private.GET("/ws/notifications", UpgradeToWebsocket(NotificationHandler))
	// AI理财宠物实时互动
	Private.GET("/ws/pet", UpgradeToWebsocket(PetInteractionHandler))
	// 好友PK实时排行榜
	Private.GET("/ws/ranking", UpgradeToWebsocket(RankingHandler))
}

// 示例 handler
func EchoHandler(conn *websocket.Conn, ctx *gin.Context) error {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		err = conn.WriteMessage(mt, message)
		if err != nil {
			return err
		}
	}
}

// 社交通知处理器
func NotificationHandler(conn *websocket.Conn, ctx *gin.Context) error {
	for {
		// 伪造实时通知处理逻辑
		mt, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		// 模拟处理通知消息
		response := map[string]interface{}{
			"type": "notification",
			"data": "Friend invitation received",
		}
		err = conn.WriteJSON(response)
		if err != nil {
			return err
		}
	}
}

// AI宠物互动处理器
func PetInteractionHandler(conn *websocket.Conn, ctx *gin.Context) error {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		// 模拟AI宠物响应
		response := map[string]interface{}{
			"type": "pet_response",
			"emotion": "happy",
			"message": "Great job saving money today!",
		}
		err = conn.WriteJSON(response)
		if err != nil {
			return err
		}
	}
}

// 排行榜实时更新处理器
func RankingHandler(conn *websocket.Conn, ctx *gin.Context) error {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		// 模拟排行榜更新
		response := map[string]interface{}{
			"type": "ranking_update",
			"rankings": []map[string]interface{}{
				{"username": "user1", "score": 1500},
				{"username": "user2", "score": 1200},
			},
		}
		err = conn.WriteJSON(response)
		if err != nil {
			return err
		}
	}
}
