package repository

import (
	"chat-management-service/models"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"time"
)

type Neo4jChatRepository struct {
	Driver neo4j.Driver
}

func NewNeo4jChatRepository(driver neo4j.Driver) *Neo4jChatRepository {
	return &Neo4jChatRepository{
		Driver: driver,
	}
}

func (repo *Neo4jChatRepository) CreateChat(chat models.ChatNode) (string, error) {
	session, err := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return "", fmt.Errorf("error creating session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	cypherQuery := `
		CREATE (c:Chat {
			dateCreated: $dateCreated,
			isActive: $isActive
		})
		RETURN c.elementId
	`

	result, err := session.Run(cypherQuery, map[string]interface{}{
		"dateCreated": chat.DateCreated,
		"isActive":    chat.IsActive,
	})
	if err != nil {
		return "", fmt.Errorf("error executing query: %v", err)
	}

	if result.Next() {
		return result.Record().Values()[0].(string), nil
	}
	return "", fmt.Errorf("no record returned")
}

func (repo *Neo4jChatRepository) SetChatToPerson(chatPerson models.ChatPerson) (string, error) {
	session, err := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return "", fmt.Errorf("error creating session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	cypherQuery := `
		MATCH (p:Person), (c:Chat)
		WHERE elementId(p) = $personElementId
		AND elementId(c) = $chatElementId
		CREATE (p)-[:PARTICIPATES_IN]->(c)
		RETURN p, c
	`

	result, err := session.Run(cypherQuery, map[string]interface{}{
		"personElementId": chatPerson.PersonElementId,
		"chatElementId":   chatPerson.ChatElementId,
	})
	if err != nil {
		return "", fmt.Errorf("error setting relation: %v", err)
	}

	if !result.Next() {
		return "", fmt.Errorf("couldn't set relation")
	}
	return "Relation done successfully", nil
}

func (repo *Neo4jChatRepository) RemoveChatToPerson(chatPerson models.ChatPerson) (string, error) {
	session, err := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return "", fmt.Errorf("error creating session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	cypherQuery := `
		MATCH (p:Person), (c:Chat)
		WHERE elementId(p) = $personElementId
		AND elementId(c) = $chatElementId
		MATCH (p)-[pi:PARTICIPATES_IN]->(c)
		DELETE pi
	`

	result, err := session.Run(cypherQuery, map[string]interface{}{
		"personElementId": chatPerson.PersonElementId,
		"chatElementId":   chatPerson.ChatElementId,
	})
	if err != nil {
		return "", fmt.Errorf("error removing relation: %v", err)
	}

	if !result.Next() {
		return "", fmt.Errorf("couldn't remove relation")
	}
	return "Relation removed successfully", nil
}

func (repo *Neo4jChatRepository) GetChatsForPerson(personElementId string) ([]models.ChatNode, error) {
	session, err := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, fmt.Errorf("error creating session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	cypherQuery := `
		MATCH (p:Person)-[:PARTICIPATES_IN]->(c:Chat)
		WHERE elementId(p) = $personElementId
		RETURN c, elementId(c) AS elementId
	`

	result, err := session.Run(cypherQuery, map[string]interface{}{
		"personElementId": personElementId,
	})
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	var chats []models.ChatNode
	for result.Next() {
		record := result.Record()

		elemIDVal, found := record.Get("elementId")
		if !found {
			continue
		}
		elemID, ok := elemIDVal.(string)
		if !ok {
			continue
		}

		cVal, found := record.Get("c")
		if !found {
			continue
		}
		node, ok := cVal.(neo4j.Node)
		if !ok {
			continue
		}

		chat := models.ChatNode{
			ElementID: elemID,
		}

		props := node.Props()

		if dateCreated, exists := props["dateCreated"]; exists {
			if t, ok := dateCreated.(time.Time); ok {
				chat.DateCreated = t
			} else if s, ok := dateCreated.(string); ok {
				parsed, err := time.Parse(time.RFC3339, s)
				if err == nil {
					chat.DateCreated = parsed
				}
			}
		}

		if isActive, exists := props["isActive"]; exists {
			if b, ok := isActive.(bool); ok {
				chat.IsActive = b
			}
		}

		chats = append(chats, chat)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating result: %v", err)
	}

	return chats, nil
}

func (repo *Neo4jChatRepository) DeleteChat(chatId string) error {
	session, err := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	cypherQuery := `
		MATCH (c:Chat)
		WHERE elementId(c) = $chatId
		DETACH DELETE c
	`

	_, err = session.Run(cypherQuery, map[string]interface{}{
		"chatId": chatId,
	})
	if err != nil {
		return fmt.Errorf("error deleting chat: %v", err)
	}

	return nil
}
