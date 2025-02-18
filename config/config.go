package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewNeo4jDriver(username, password, url string) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("error connecting to Neo4j: %v", err)
	}
	return driver, nil
}

func NewMongoClient(mongoURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("error verifying MongoDB connection: %v", err)
	}

	return client, nil
}

func NewRedisClient(redisAddr string) (*redis.Client, error) {
	if redisAddr == "" {
		return nil, fmt.Errorf("redis address is empty")
	}

	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, fmt.Errorf("error parsing Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %v. RedisAddr: %s", err, redisAddr)
	}

	return client, nil
}
