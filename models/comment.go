package models

import (
	"time"
)

// Comment представляет комментарий к новости
type Comment struct {
	ID        int
	NewsID    int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment создает новый комментарий с заданными параметрами
func NewComment(newsID, userID int, username, content string) *Comment {
	now := time.Now()
	return &Comment{
		NewsID:    newsID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
} 