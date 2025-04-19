package database

import (
	"forum-app/helpers"
	"forum-app/models"

	"strings"
	"time"
)

func (db *Connection) SetPost(title, content, author, categories string) error {
	// Sanitize input
	cleanTitle, cleanContent, err := helpers.SanitizePost(title, content)
	if err != nil {
		return err
	}

	query := `INSERT INTO post(title, categories, content, author, time)
				VALUES(?, ?, ?, ?, ?)`

	_, err = db.DB.Exec(query, cleanTitle, categories, cleanContent, author, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func (db *Connection) GetTotalPostCount(filter string, user *models.Users) (int, error) {
	query := `SELECT COUNT(*) FROM post`
	var args []interface{}

	// Apply filter if provided
	if filter != "" {
		if filter == "Created" {
			query += ` WHERE author = ?`
			args = append(args, user.ID)
		} else {
			query += ` WHERE categories LIKE ?`
			args = append(args, "%"+filter+"%")
		}
	}

	var count int
	err := db.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Connection) GetPostsForHome(page int, filter string, user *models.Users) ([]models.Post, error) {
	const pageSize = 10
	offset := (page - 1) * pageSize
	var args []interface{}
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes 
              FROM post p 
              JOIN user u ON p.author = u.id`
	if filter != "" && filter != "Created" && filter != "Liked" {
		query += ` WHERE p.categories LIKE ?`
		args = append(args, "%"+filter+"%")
	}
	if filter == "Created" {
		query += ` WHERE p.author = ?`
		args = append(args, user.ID)
	}
	query += ` ORDER BY p.time DESC 
			  LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	var user_id int
	for rows.Next() {
		var post models.Post
		var categories string
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&categories,
			&post.Content,
			&user_id,
			&post.Time,
			&post.Upvotes,
			&post.Downvotes,
		)
		if err != nil {
			return nil, err
		}
		post.Author, err = db.GetUserById(user_id)
		if err != nil {
			return nil, err
		}
		post.Categories = strings.Split(categories, ",")
		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Connection) GetPostByID(id int) (models.Post, error) {
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes 
              FROM post p 
              JOIN user u ON p.author = u.id 
              WHERE p.id = ?`
	var post models.Post

	var categories string
	var user_id int
	err := db.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&categories,
		&post.Content,
		&user_id,
		&post.Time,
		&post.Upvotes,
		&post.Downvotes,
	)
	if err != nil {
		return post, err
	}
	post.Author, err = db.GetUserById(user_id)
	if err != nil {
		return post, err
	}

	post.Categories = strings.Split(categories, ",")

	// Get comments for the post
	commentsQuery := `SELECT c.content, c.author, c.time, c.upvotes, c.downvotes

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
		var user_id int
		err := rows.Scan(&comment.Content, &user_id, &comment.Time, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return post, err
		}
		comment.Author, err = db.GetUserById(user_id)
		if err != nil {
			return post, err
		}
		post.Comments = append(post.Comments, comment)
	}

	return post, nil
}
