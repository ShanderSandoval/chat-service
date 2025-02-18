package repository

import (
	"chat-management-service/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisChatRepository struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisChatRepository(client *redis.Client) *RedisChatRepository {
	return &RedisChatRepository{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (repo *RedisChatRepository) chatKey(chatId string) string {
	return fmt.Sprintf("chat:%s", chatId)
}

func (repo *RedisChatRepository) CreateChat(chat models.ChatVolatile) (string, error) {
	if chat.DateCreated.IsZero() {
		chat.DateCreated = time.Now()
	}
	key := repo.chatKey(chat.Id)
	data, err := json.Marshal(chat)
	if err != nil {
		return "", fmt.Errorf("error marshalling chat: %v", err)
	}
	err = repo.Client.Set(repo.Ctx, key, data, 0).Err()
	if err != nil {
		return "", fmt.Errorf("error storing chat in redis: %v", err)
	}
	return chat.Id, nil
}

func (repo *RedisChatRepository) GetChat(chatId string) (*models.ChatVolatile, error) {
	key := repo.chatKey(chatId)
	data, err := repo.Client.Get(repo.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting chat from redis: %v", err)
	}
	var chat models.ChatVolatile
	err = json.Unmarshal([]byte(data), &chat)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling chat: %v", err)
	}
	return &chat, nil
}

func (repo *RedisChatRepository) UpdateChat(chat models.ChatVolatile) error {
	key := repo.chatKey(chat.Id)
	data, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("error marshalling chat: %v", err)
	}
	err = repo.Client.Set(repo.Ctx, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("error updating chat in redis: %v", err)
	}
	return nil
}

func (repo *RedisChatRepository) AddMessageToChat(chatId string, message models.ChatMessage) error {
	chat, err := repo.GetChat(chatId)
	if err != nil {
		return err
	}
	chat.Messages = append(chat.Messages, message)
	return repo.UpdateChat(*chat)
}

func (repo *RedisChatRepository) DeleteChat(chatId string) error {
	key := repo.chatKey(chatId)
	err := repo.Client.Del(repo.Ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting chat from redis: %v", err)
	}
	return nil
}

func (repo *RedisChatRepository) FindChatById(chatId string) (*models.ChatVolatile, error) {
	return repo.GetChat(chatId)
}
