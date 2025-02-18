package models

import "time"

type ChatMessage struct {
	Date            time.Time `json:"date"`
	PersonElementId string    `json:"personElementId"`
	Body            string    `json:"body"`
}
