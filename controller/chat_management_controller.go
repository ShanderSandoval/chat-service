package controller

import (
	"chat-management-service/models"
	"chat-management-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	ChatService *service.ChatService
}

func NewChatController(chatService *service.ChatService) *ChatController {
	return &ChatController{
		ChatService: chatService,
	}
}

func (cc *ChatController) RegisterRoutes(router *gin.Engine) {
	router.POST("/chatService", cc.CreateChat)
	router.POST("/chatService/:id/message", cc.AddMessage)
	router.PUT("/chatService/:id/sync", cc.SyncMessages)
	router.DELETE("/chatService/:id/redis", cc.DeleteChatFromRedis)
	router.POST("/chatService/person", cc.AddPersonToChat)
	router.DELETE("/chatService/person", cc.RemovePersonFromChat)
	router.GET("/chatService/person/:personElementId", cc.GetChatsForPerson)
	router.DELETE("/chatService/:id", cc.DeleteChat)
}

func (cc *ChatController) CreateChat(c *gin.Context) {
	var chatNode models.ChatNode
	if err := c.ShouldBindJSON(&chatNode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	elementId, err := cc.ChatService.CreateChat(chatNode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chatId": elementId})
}

func (cc *ChatController) AddMessage(c *gin.Context) {
	chatId := c.Param("id")
	var message models.ChatMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.ChatService.AddMessage(chatId, message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message added successfully"})
}

func (cc *ChatController) SyncMessages(c *gin.Context) {
	chatId := c.Param("id")
	if err := cc.ChatService.SyncMessages(chatId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages synced successfully"})
}

func (cc *ChatController) DeleteChatFromRedis(c *gin.Context) {
	chatId := c.Param("id")
	if err := cc.ChatService.DeleteChatFromRedis(chatId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted from Redis successfully"})
}

func (cc *ChatController) AddPersonToChat(c *gin.Context) {
	var chatPerson models.ChatPerson
	if err := c.ShouldBindJSON(&chatPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := cc.ChatService.AddPersonToChat(chatPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

func (cc *ChatController) RemovePersonFromChat(c *gin.Context) {
	var chatPerson models.ChatPerson
	if err := c.ShouldBindJSON(&chatPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := cc.ChatService.RemovePersonFromChat(chatPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

func (cc *ChatController) GetChatsForPerson(c *gin.Context) {
	personElementId := c.Param("personElementId")
	chats, err := cc.ChatService.GetChatsForPerson(personElementId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (cc *ChatController) DeleteChat(c *gin.Context) {
	chatId := c.Param("id")
	if err := cc.ChatService.DeleteChat(chatId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
