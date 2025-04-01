package database

import (
	"forum-app/models"
	"time"
)

func (db *Connection) SetPost(title, category, content, author string) error {
	query := `INSERT INTO post(title, category, content, author, time)
				VALUES(?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, title, category, content, author, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func (db *Connection) GetPostByID(id int) (models.Post, error) {
	query := `SELECT p.id, p.title, p.category, p.content, u.username, p.time, p.likes 
              FROM post p 
              JOIN user u ON p.author = u.id 
              WHERE p.id = ?`
	var post models.Post

	err := db.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Category,
		&post.Content,
		&post.Author,
		&post.Time,
		&post.Likes,
	)
	if err != nil {
		return post, err
	}

	// Get comments for the post
	commentsQuery := `SELECT c.content, u.username, c.time, c.likes 
                     FROM comment c 
                     JOIN user u ON c.author = u.id 
                     WHERE c.post_id = ?`

	rows, err := db.DB.Query(commentsQuery, id)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.Content, &comment.Author, &comment.Time, &comment.Likes)
		if err != nil {
			return post, err
		}
		post.Comments = append(post.Comments, comment)
	}

	return post, nil
}
