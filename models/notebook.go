package models

import "time"

type Notebook struct {
	Id           string
	Title        string `json:"title"`
	Content      string `json:"content"`
	Topic        string `json:"topic"`
	Owner        string // user_id
	Deleted      bool
	LastModified time.Time
	CreatedAt    time.Time
}
