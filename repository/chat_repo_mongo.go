package repository

import (
	"chat-management-service/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoChatRepository struct {
	Collection *mongo.Collection
}

func NewMongoChatRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRepository{Collection: collection}
}

func (repo *MongoChatRepository) CreateChat(chat models.ChatCollection) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if chat.DateCreated.IsZero() {
		chat.DateCreated = time.Now()
	}

	_, err := repo.Collection.InsertOne(ctx, chat)
	if err != nil {
		return "", fmt.Errorf("error creating chat: %v", err)
	}
	return chat.Id, nil
}

func (repo *MongoChatRepository) AddMessageToChat(chatId string, message models.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": chatId}
	update := bson.M{"$push": bson.M{"messages": message}}
	_, err := repo.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error adding message to chat: %v", err)
	}
	return nil
}

func (repo *MongoChatRepository) DeactivateChat(chatId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": chatId}
	update := bson.M{"$set": bson.M{"isActive": false}}
	_, err := repo.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error deactivating chat: %v", err)
	}
	return nil
}

func (repo *MongoChatRepository) DeleteChat(chatId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": chatId}
	_, err := repo.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting chat: %v", err)
	}
	return nil
}

func (repo *MongoChatRepository) FindChatById(chatId string) (*models.ChatCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": chatId}
	var chat models.ChatCollection
	err := repo.Collection.FindOne(ctx, filter).Decode(&chat)
	if err != nil {
		return nil, fmt.Errorf("error finding chat by id: %v", err)
	}
	return &chat, nil
}

func (repo *MongoChatRepository) UpdateChatMessages(chatId string, messages []models.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": chatId}
	update := bson.M{"$set": bson.M{"messages": messages}}
	_, err := repo.Collection.UpdateOne(ctx, filter, update)
	return err
}
