package models

import (
	"time"
)

// News представляет новость в системе
type News struct {
	ID        int
	Title     string
	Content   string
	AuthorID  int
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewNews создает новую новость с заданными параметрами
func NewNews(title, content string, authorID int, author string) *News {
	now := time.Now()
	return &News{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
	}
} 