package database

import (
	"forum-app/models"
	"time"
)

func (db *Connection) SetPost(title, category, content, author string) error {
	query := `INSERT INTO post(title, category, content, author)
				VALUES(?, ?, ?, ?)`

	_, err := db.DB.Exec(query, title, category, content, author, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func (db *Connection) GetPostByID(id int) (models.Post, error) {
	query := `SELECT * FROM user WHERE id = ? `
	var post models.Post

	err := db.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Category, &post.Content, &post.Author, &post.Time, &post.Likes, &post.Comments)

	return post, err
}
