package main
import (
	"github.com/gin-gonic/gin"
	"porter/api/src/handlers/websocket"
	"porter/api/src/handlers/system"
)

func main() {
	r := gin.New()
	r.GET("/ping", system.Ping)
	r.GET("/ws", websocket.Handler)
	r.Run(":1111")
}