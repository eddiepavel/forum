package database

import (
	"database/sql"
	"forum-app/models"
	"strings"
	"time"
)

func (db *Connection) buildHomeQuery(filter string, user *models.Users) (string, []interface{}) {
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes, 
                 (SELECT COUNT(*) FROM comment c WHERE c.post_id = p.id) AS comment_count
          FROM post p 
          JOIN user u ON p.author = u.id`
	var args []interface{}

	switch filter {
	case "Created":
		query += ` WHERE p.author = ?`
		args = append(args, user.ID)
	case "Liked":
		query += ` WHERE p.id IN (SELECT post_id FROM votes WHERE user_id = ? AND vote_type = 'upvote')`
		args = append(args, user.ID)
	default:
		if filter != "" {
			query += ` WHERE p.categories LIKE ?`
			args = append(args, "%"+filter+"%")
		}
	}

	return query, args
}

func (db *Connection) scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post

	for rows.Next() {
		post, err := db.scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Connection) scanPostRow(rows *sql.Rows) (models.Post, error) {
	var post models.Post
	var categories string
	var userID int
	var timeRaw time.Time

	err := rows.Scan(
		&post.ID,
		&post.Title,
		&categories,
		&post.Content,
		&userID,
		&timeRaw,
		&post.Upvotes,
		&post.Downvotes,
		&post.CommentCount,
	)
	if err != nil {
		return post, err
	}

	post.Author, err = db.GetUserById(userID)
	if err != nil {
		return post, err
	}

	post.Time = timeRaw.Format("2006-01-02 15:04:05")
	post.Categories = strings.Split(categories, ",")

	return post, nil
}

func (db *Connection) fetchPostByID(id int) (models.Post, error) {
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes, p.vote_count 
              FROM post p 
              JOIN user u ON p.author = u.id 
              WHERE p.id = ?`
	var post models.Post
	var categories string
	var userID int
	var timeRaw time.Time

	err := db.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&categories,
		&post.Content,
		&userID,
		&timeRaw,
		&post.Upvotes,
		&post.Downvotes,
		&post.VoteCount,
	)
	if err != nil {
		return post, err
	}

	post.Author, err = db.GetUserById(userID)
	if err != nil {
		return post, err
	}

	post.Time = timeRaw.Format("2006-01-02 15:04:05")
	post.Categories = strings.Split(categories, ",")

	return post, nil
}

func (db *Connection) fetchPostComments(postID, userID int) ([]models.Comment, error) {
	query := `SELECT c.id, c.content, c.author, c.time, c.upvotes, c.downvotes, c.vote_count FROM comment c WHERE c.post_id = ?`
	rows, err := db.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		comment, err := db.scanCommentRow(rows, userID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *Connection) scanCommentRow(rows *sql.Rows, userID int) (models.Comment, error) {
	var comment models.Comment
	var timeRaw time.Time
	var authorID int

	err := rows.Scan(
		&comment.ID,
		&comment.Content,
		&authorID,
		&timeRaw,
		&comment.Upvotes,
		&comment.Downvotes,
		&comment.VoteCount,
	)
	if err != nil {
		return comment, err
	}

	comment.Author, err = db.GetUserById(authorID)
	if err != nil {
		return comment, err
	}

	comment.UserVote, err = db.GetUserVote(userID, 0, comment.ID)
	if err != nil && err != sql.ErrNoRows {
		return comment, err
	}

	comment.Time = timeRaw.Format("2006-01-02 15:04:05")

	return comment, nil
}

func (db *Connection) updateVoteCounts(postID, commentID int) {
	if postID != 0 {
		db.DB.Exec(`UPDATE post SET 
            upvotes = (SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'upvote'),
            downvotes = (SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'downvote'),
            vote_count = (upvotes - downvotes)
            WHERE id = ?`, postID, postID, postID)
	} else if commentID != 0 {
		db.DB.Exec(`UPDATE comment SET 
            upvotes = (SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'upvote'),
            downvotes = (SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'downvote'),
            vote_count = (upvotes - downvotes)
            WHERE id = ?`, commentID, commentID, commentID)
	}
}
