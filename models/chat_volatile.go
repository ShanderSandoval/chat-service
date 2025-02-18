package models

import "time"

type ChatVolatile struct {
	Id          string        `json:"id"`
	DateCreated time.Time     `json:"dateCreated"`
	IsActive    bool          `json:"isActive"`
	Messages    []ChatMessage `json:"messages"`
}
