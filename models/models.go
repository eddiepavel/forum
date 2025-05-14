package models

import (
	"forum-app/session"
	"html/template"
	"time"
)

type Users struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Auth      string    `json:"auth"`
	Picture   string    `json:"picture"`
	Is_Admin  int       `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
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
