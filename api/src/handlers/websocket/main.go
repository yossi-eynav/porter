package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"time"
	"porter/api/src/lib/github"
	"porter/api/src/handlers/websocket/message"
	"fmt"
	"porter/api/src/handlers/websocket/events"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReceivedSocketEvent struct {
	Type string `json:"type"`
}

func Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
§
		var event ReceivedSocketEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handlerSockerEvent(event, conn)
	}
}

func handlerSockerEvent(event ReceivedSocketEvent, conn *websocket.Conn) {
	switch event.Type {
		case events.EventGetUsedPorts:
			messagesChannel := make(chan message.Message)
			repositoriesChannel := make(chan github.Repository)
			go usedPortsHandler(repositoriesChannel, messagesChannel)
			go func(){
				for  {
					select {
						case repo, ok := <-repositoriesChannel:
							if !ok {
								msg, _ :=json.Marshal(message.Message{
									Body: "Done.",
									Color: "green",
									Timestamp: time.Now(),
								})
								conn.WriteMessage(websocket.TextMessage, msg)
								return
							}
							msg, _ :=json.Marshal(repo)
							conn.WriteMessage(websocket.TextMessage, msg)
						case msg, _ := <-messagesChannel:
							strigifiedMsg, _ :=json.Marshal(msg)
							conn.WriteMessage(websocket.TextMessage, strigifiedMsg)
					}
				}
			}()
	}
}