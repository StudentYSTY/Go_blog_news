package models

import (
	"time"
)

// User представляет пользователя в системе
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	IsBlocked bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser создает нового пользователя с заданными параметрами
func NewUser(username, password, email string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Password:  password,
		Email:     email,
		IsBlocked: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
} 