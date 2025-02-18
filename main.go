package main

import (
	"chat-management-service/config"
	"chat-management-service/controller"
	"chat-management-service/repository"
	"chat-management-service/service"
	"chat-management-service/ws"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hoshsadiq/godotenv"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	redisURI := os.Getenv("REDIS_URI")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	neo4jDriver, err := config.NewNeo4jDriver(neo4jUsername, neo4jPassword, neo4jURI)
	if err != nil {
		log.Fatalf("Error creating Neo4j driver: %v", err)
	}
	defer func(neo4jDriver neo4j.Driver) {
		err := neo4jDriver.Close()
		if err != nil {

		}
	}(neo4jDriver)

	mongoClient, err := config.NewMongoClient(mongoURI)
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {

		}
	}(mongoClient, nil)

	redisClient, err := config.NewRedisClient(redisURI)
	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {

		}
	}(redisClient)

	neoRepo := repository.NewNeo4jChatRepository(neo4jDriver)
	mongoRepo := repository.NewMongoChatRepository(mongoClient, mongoDatabase, mongoCollection)
	redisRepo := repository.NewRedisChatRepository(redisClient)

	chatService := service.NewChatService(neoRepo, mongoRepo, redisRepo)

	r := gin.Default()

	chatController := controller.NewChatController(chatService)

	chatController.RegisterRoutes(r)

	r.GET("/ws", func(c *gin.Context) {
		ws.HandleConnectionsGin(c, chatService)
	})

	err = r.Run(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return
	}
}
