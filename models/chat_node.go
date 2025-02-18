package models

import "time"

type ChatNode struct {
	ElementID   string    `json:"elementId"`
	DateCreated time.Time `json:"dateCreated"`
	IsActive    bool      `json:"isActive"`
}
