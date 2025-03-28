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

type Post struct {
	ID       int
	Title    string
	Category string
	Content  string
	Author   string
	Time     time.Time
	Likes    int
	Comments []Comment
}

type Comment struct {
	Content string
	Author  string
	Time    time.Time
	Likes   int
}
