package ws

import (
	"chat-management-service/models"
	"chat-management-service/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request, chatService *service.ChatService) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing WebSocket:", err)
		}
	}(conn)

	chatId := r.URL.Query().Get("chatId")

	if chatId == "" {
		log.Println("No chatId provided")
		return
	}

	err = chatService.SyncMessages(chatId)
	if err != nil {
		log.Println("Error syncing messages:", err)
		return
	}

	redisChat, err := chatService.GetChatsForPerson(chatId)
	if err != nil {
		log.Println("Error retrieving chat from Redis:", err)
		return
	}

	err = conn.WriteJSON(redisChat)
	if err != nil {
		log.Println("Error sending messages to client:", err)
		return
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		if msgType == websocket.TextMessage {
			message := models.ChatMessage{
				Body: string(msg),
			}

			err = chatService.AddMessage(chatId, message)
			if err != nil {
				log.Println("Error saving message to Redis:", err)
				break
			}

			err = conn.WriteMessage(msgType, msg)
			if err != nil {
				log.Println("Error sending message to client:", err)
				break
			}
		}
	}
}

func HandleConnectionsGin(c *gin.Context, chatService *service.ChatService) {
	HandleConnections(c.Writer, c.Request, chatService)
}
