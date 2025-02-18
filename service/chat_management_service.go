package service

import (
	"chat-management-service/models"
	"chat-management-service/repository"
	"fmt"
	"sort"
)

type ChatService struct {
	neoRepo   *repository.Neo4jChatRepository
	mongoRepo *repository.MongoChatRepository
	redisRepo *repository.RedisChatRepository
}

func NewChatService(neoRepo *repository.Neo4jChatRepository, mongoRepo *repository.MongoChatRepository, redisRepo *repository.RedisChatRepository) *ChatService {
	return &ChatService{
		neoRepo:   neoRepo,
		mongoRepo: mongoRepo,
		redisRepo: redisRepo,
	}
}

func (s *ChatService) CreateChat(chatNeo models.ChatNode) (string, error) {
	elementId, err := s.neoRepo.CreateChat(chatNeo)
	if err != nil {
		return "", fmt.Errorf("failed to create chat in Neo4j: %v", err)
	}

	chatMongo := models.ChatCollection{
		Id:          elementId,
		DateCreated: chatNeo.DateCreated,
		IsActive:    chatNeo.IsActive,
		Messages:    []models.ChatMessage{},
	}

	chatRedis := models.ChatVolatile{
		Id:          elementId,
		DateCreated: chatNeo.DateCreated,
		IsActive:    chatNeo.IsActive,
		Messages:    []models.ChatMessage{},
	}

	if _, err := s.mongoRepo.CreateChat(chatMongo); err != nil {
		return "", fmt.Errorf("failed to create chat in MongoDB: %v", err)
	}

	if _, err := s.redisRepo.CreateChat(chatRedis); err != nil {
		return "", fmt.Errorf("failed to create chat in Redis: %v", err)
	}

	return elementId, nil
}

func (s *ChatService) AddMessage(chatId string, message models.ChatMessage) error {
	return s.redisRepo.AddMessageToChat(chatId, message)
}

func (s *ChatService) SyncMessages(chatId string) error {
	redisChat, err := s.redisRepo.FindChatById(chatId)
	if err != nil {
		return fmt.Errorf("failed to get chat from Redis: %v", err)
	}

	mongoChat, err := s.mongoRepo.FindChatById(chatId)
	if err != nil {
		return fmt.Errorf("failed to get chat from MongoDB: %v", err)
	}

	merged := mergeMessages(mongoChat.Messages, redisChat.Messages)

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Date.Before(merged[j].Date)
	})

	if err := s.mongoRepo.UpdateChatMessages(chatId, merged); err != nil {
		return fmt.Errorf("failed to update chat messages in MongoDB: %v", err)
	}

	return nil
}

func (s *ChatService) DeleteChatFromRedis(chatId string) error {
	return s.redisRepo.DeleteChat(chatId)
}

func (s *ChatService) AddPersonToChat(chatPerson models.ChatPerson) (string, error) {
	return s.neoRepo.SetChatToPerson(chatPerson)
}

func (s *ChatService) RemovePersonFromChat(chatPerson models.ChatPerson) (string, error) {
	return s.neoRepo.RemoveChatToPerson(chatPerson)
}

func (s *ChatService) GetChatsForPerson(personElementId string) (interface{}, error) {
	chats, err := s.neoRepo.GetChatsForPerson(personElementId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats for person in Neo4j: %v", err)
	}
	return chats, nil
}

func (s *ChatService) DeleteChat(chatId string) error {
	if err := s.neoRepo.DeleteChat(chatId); err != nil {
		return fmt.Errorf("failed to delete chat in Neo4j: %v", err)
	}
	if err := s.mongoRepo.DeleteChat(chatId); err != nil {
		return fmt.Errorf("failed to delete chat in MongoDB: %v", err)
	}
	if err := s.redisRepo.DeleteChat(chatId); err != nil {
		return fmt.Errorf("failed to delete chat in Redis: %v", err)
	}
	return nil
}

func mergeMessages(existing, incoming []models.ChatMessage) []models.ChatMessage {
	messageMap := make(map[string]models.ChatMessage)
	for _, msg := range existing {
		key := fmt.Sprintf("%d-%s", msg.Date.UnixNano(), msg.Body)
		messageMap[key] = msg
	}
	for _, msg := range incoming {
		key := fmt.Sprintf("%d-%s", msg.Date.UnixNano(), msg.Body)
		if _, found := messageMap[key]; !found {
			messageMap[key] = msg
		}
	}
	merged := make([]models.ChatMessage, 0, len(messageMap))
	for _, msg := range messageMap {
		merged = append(merged, msg)
	}
	return merged
}
