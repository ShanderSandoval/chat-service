package models

import "time"

type ChatCollection struct {
	Id          string        `json:"id"`
	DateCreated time.Time     `json:"dateCreated"`
	IsActive    bool          `json:"isActive"`
	Messages    []ChatMessage `json:"messages"`
}
