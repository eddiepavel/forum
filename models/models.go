package models

import "time"

type Users struct {
	ID        int
	Email     string
	Username  string
	Password  string
	Is_Admin  int
	CreatedAt time.Time
}

type Session struct {
	ID        int
	Token     string
	ExpiresAt time.Time
	UserId    int
}
