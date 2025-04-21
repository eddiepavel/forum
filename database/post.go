package database

import (
	"database/sql"
	"errors"
	"fmt"
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
		} else if filter == "Liked" {
			query += ` WHERE id IN (SELECT post_id FROM votes WHERE user_id = ? AND vote_type = 'upvote')`
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
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes, 
                 (SELECT COUNT(*) FROM comment c WHERE c.post_id = p.id) AS comment_count
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
	if filter == "Liked" {
		query += ` WHERE p.id IN (SELECT post_id FROM votes WHERE user_id = ? AND vote_type = 'upvote')`
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
	var timeRaw time.Time
	for rows.Next() {
		var post models.Post
		var categories string
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&categories,
			&post.Content,
			&user_id,
			&timeRaw,
			&post.Upvotes,
			&post.Downvotes,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		post.Author, err = db.GetUserById(user_id)
		if err != nil {
			return nil, err
		}

		post.Time = timeRaw.Format("2006-01-02 15:04:05")
		post.Categories = strings.Split(categories, ",")
		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Connection) GetPostByID(id int, userID int) (models.Post, error) {
	query := `SELECT p.id, p.title, p.categories, p.content, p.author, p.time, p.upvotes, p.downvotes, p.vote_count 
              FROM post p 
              JOIN user u ON p.author = u.id 
              WHERE p.id = ?`
	var post models.Post

	var categories string
	var user_id int
	var timeRaw time.Time
	err := db.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&categories,
		&post.Content,
		&user_id,
		&timeRaw,
		&post.Upvotes,
		&post.Downvotes,
		&post.VoteCount,
	)
	if err != nil {
		return post, err
	}
	post.Author, err = db.GetUserById(user_id)
	if err != nil {
		return post, err
	}
	post.Time = timeRaw.Format("2006-01-02 15:04:05")
	post.Categories = strings.Split(categories, ",")

	var voteType string
	err = db.DB.QueryRow(`SELECT vote_type FROM votes WHERE user_id = ? AND post_id = ?`, userID, id).Scan(&voteType)
	if err == nil {
		post.UserVote = voteType
	} else {
		post.UserVote = "none"
	}

	// Get comments and their vote states
	commentsQuery := `SELECT c.id, c.content, c.author, c.time, c.upvotes, c.downvotes, c.vote_count FROM comment c WHERE c.post_id = ?`
	rows, err := db.DB.Query(commentsQuery, id)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var timeRaw time.Time
		var user_id int
		err := rows.Scan(&comment.ID, &comment.Content, &user_id, &timeRaw, &comment.Upvotes, &comment.Downvotes, &comment.VoteCount)
		if err != nil {
			return post, err
		}
		comment.Author, err = db.GetUserById(user_id)
		if err != nil {
			return post, err
		}
		err = db.DB.QueryRow(`SELECT vote_type FROM votes WHERE user_id = ? AND comment_id = ?`, userID, comment.ID).Scan(&voteType)
		if err == nil {
			comment.UserVote = voteType
		} else {
			comment.UserVote = "none"
		}
		comment.Time = timeRaw.Format("2006-01-02 15:04:05")
		post.Comments = append(post.Comments, comment)
	}

	return post, nil
}

func (db *Connection) SetVote(userID, postID, commentID int, voteType string) error {
	// Check if a vote already exists
	var existingVote string
	query := `SELECT vote_type FROM votes WHERE user_id = ? AND post_id = ? AND comment_id = ?`
	err := db.DB.QueryRow(query, userID, postID, commentID).Scan(&existingVote)

	if err == nil {
		// If the vote exists and is the same, remove it
		if existingVote == voteType {
			_, err := db.DB.Exec(`DELETE FROM votes WHERE user_id = ? AND post_id = ? AND comment_id = ?`, userID, postID, commentID)
			if err == nil {
				db.updateVoteCounts(postID, commentID)
			}
			return err
		}

		// If the vote exists but is different, update it
		_, err := db.DB.Exec(`UPDATE votes SET vote_type = ? WHERE user_id = ? AND post_id = ? AND comment_id = ?`, voteType, userID, postID, commentID)
		if err == nil {
			db.updateVoteCounts(postID, commentID)
		}
		return err
	}

	// If no vote exists, insert a new one
	_, err = db.DB.Exec(`INSERT INTO votes (user_id, post_id, comment_id, vote_type) VALUES (?, ?, ?, ?)`, userID, postID, commentID, voteType)
	if err == nil {
		db.updateVoteCounts(postID, commentID)
	}
	return err
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

func (db *Connection) GetPostVoteCounts(postID int) (int, int) {
	var upvotes, downvotes int
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'upvote'`, postID).Scan(&upvotes)
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'downvote'`, postID).Scan(&downvotes)
	return upvotes, downvotes
}

func (db *Connection) GetCommentVoteCounts(commentID int) (int, int) {
	var upvotes, downvotes int
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'upvote'`, commentID).Scan(&upvotes)
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'downvote'`, commentID).Scan(&downvotes)
	return upvotes, downvotes
}

func (db *Connection) GetUserVote(userID, postID, commentID int) string {
	var voteType string
	err := db.DB.QueryRow(`SELECT vote_type FROM votes WHERE user_id = ? AND post_id = ? AND comment_id = ?`, userID, postID, commentID).Scan(&voteType)
	if err != nil {
		return "none" // Default to "none" if no vote exists
	}
	return voteType
}

func (db *Connection) DeletePost(postID, userID int) error {
	// Check if the user is the author of the post
	var authorID int
	err := db.DB.QueryRow(`SELECT author FROM post WHERE id = ?`, postID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("post not found")
		}
		return err
	}

	if authorID != userID {
		return errors.New("403 Forbidden: You are not the author of this post")
	}

	// Proceed with deletion
	query := `DELETE FROM post WHERE id = ? AND author = ?`
	result, err := db.DB.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}
