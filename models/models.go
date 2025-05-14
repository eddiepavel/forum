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
	Auth      string
	Picture   string
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
	Image        string
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

type Notification struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postID"`
	Target  Users  `json:"target"`
	Actor   Users  `json:"actor"`
	Time    string `json:"time"`
	Type    string `json:"type"`
	Read    bool   `json:"read"`
	Content string `json:"content"`
}

type PageData struct {
	Data     map[string]interface{}
	User     *Users
	Session  *session.Session
	Source   string
	Redirect string
}
