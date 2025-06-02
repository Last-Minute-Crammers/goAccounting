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
	// 可继续注册其他 websocket 路由
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
