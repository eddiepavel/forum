package models

import (
	"forum-app/session"
	"html/template"
	"time"
)

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

type Post struct {
	ID        int
	Title     string
	Category  string
	Content   template.HTML
	Author    Users
	Time      string
	Upvotes   int
	Downvotes int
	Comments  []Comment
}

type Comment struct {
	PostID    int
	Content   template.HTML
	Author    Users
	Time      string
	Upvotes   int
	Downvotes int
}

type PageData struct {
	Data     map[string]interface{}
	User     *Users
	Session  *session.Session
	Source   string
	Redirect string
}
