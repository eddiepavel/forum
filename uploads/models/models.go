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
	ID           int
	Title        string
	Categories   []string
	Content      template.HTML
	Author       Users
	Time         string
	Upvotes      int
	Downvotes    int
	VoteCount    int
	CommentCount int
	Comments     []Comment
	UserVote     string
}

type Comment struct {
	ID        int
	PostID    int
	Content   template.HTML
	Author    Users
	Time      string
	Upvotes   int
	Downvotes int
	VoteCount int
	UserVote  string
}

type PageData struct {
	Data     map[string]interface{}
	User     *Users
	Session  *session.Session
	Source   string
	Redirect string
}
