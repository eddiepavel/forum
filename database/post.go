package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-app/helpers"
	"forum-app/models"

	"time"
)

// SetPost inserts a new post into the database with the given title, content, author, and categories.
func (db *Connection) SetPost(title, content, author, categories, imagePath string) error {
	// Sanitize input
	cleanTitle, cleanContent, err := helpers.SanitizePost(title, content)
	if err != nil {
		return err
	}

	query := `INSERT INTO post(title, categories, content, author, time, image_path)
          VALUES(?, ?, ?, ?, ?, ?)`

	_, err = db.DB.Exec(query, cleanTitle, categories, cleanContent, author, time.Now().Format("2006-01-02 15:04:05"), imagePath)
	return err
}

// GetTotalPostCount retrieves the total number of posts based on the provided filter and user context.
func (db *Connection) GetTotalPostCount(filter string, user *models.Users) (int, error) {
	if filter == "Liked" || filter == "Created" {
		if user == nil {
			return 0, errors.New("user not logged in")
		}
	}
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

// GetPostsForHome retrieves a paginated list of posts for the home page based on the filter and user context.
func (db *Connection) GetPostsForHome(page int, filter string, user *models.Users) ([]models.Post, error) {
	const pageSize = 10
	offset := (page - 1) * pageSize

	query, args := db.buildHomeQuery(filter, user)
	query += ` ORDER BY p.time DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanPosts(rows)
}

func (db *Connection) GetUserActivity(userID int) ([]models.UserPostActivity, error) {
	var activity []models.UserPostActivity

	// 1. User's created posts
	qCreated := `
		SELECT 
			p.id, p.title, p.author, p.time, SUBSTR(p.content, 1, 200) as content
		FROM post p
		WHERE p.author = ?
`
	rows, err := db.DB.Query(qCreated, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.UserPostActivity
		var ts time.Time
		if err := rows.Scan(&a.PostID, &a.PostTitle, &a.Author, &ts, &a.ContentPreview); err != nil {
			return nil, err
		}
		a.EventType = "created"
		a.PostTimestamp = ts.Format("2006-01-02 15:04:05")
		activity = append(activity, a)
	}
	rows.Close()

	// 2. Posts where user has left up/down vote
	qVoted := `
		SELECT 
			p.id, p.title, p.author, p.time, SUBSTR(p.content, 1, 200) as content,
			v.vote_type
		FROM post p
		JOIN votes v ON v.post_id = p.id
		WHERE v.user_id = ? AND v.comment_id = 0
`
	rows, err = db.DB.Query(qVoted, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.UserPostActivity
		var ts time.Time
		if err := rows.Scan(&a.PostID, &a.PostTitle, &a.Author, &ts, &a.ContentPreview, &a.VoteType); err != nil {
			return nil, err
		}
		a.EventType = "voted"
		a.PostTimestamp = ts.Format("2006-01-02 15:04:05")
		activity = append(activity, a)
	}
	rows.Close()

	// 3. Posts where the user has commented with what the user commented
	qCommented := `
		SELECT 
			p.id, p.title, p.author, p.time, SUBSTR(p.content, 1, 200) as post_content,
			SUBSTR(c.content, 1, 200) as comment_content,
			c.time
		FROM post p
		JOIN comment c ON c.post_id = p.id
		WHERE c.author = ?
		`
	rows, err = db.DB.Query(qCommented, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.UserPostActivity
		var ts time.Time
		var cts string
		if err := rows.Scan(
			&a.PostID, &a.PostTitle, &a.Author, &ts, &a.ContentPreview, &a.CommentContent, &cts,
		); err != nil {
			return nil, err
		}
		a.EventType = "commented"
		a.PostTimestamp = ts.Format("2006-01-02 15:04:05")
		a.CommentTimestamp = cts
		activity = append(activity, a)
	}
	rows.Close()

	return activity, nil

}

// GetPostByID retrieves a post by its ID, including user-specific vote and comment data.
func (db *Connection) GetPostByID(id int, userID int) (models.Post, error) {
	post, err := db.fetchPostByID(id)
	if err != nil {
		return post, err
	}

	post.UserVote, err = db.GetUserVote(userID, id, 0)
	if err != nil && err != sql.ErrNoRows {
		return post, err
	}

	post.Comments, err = db.fetchPostComments(id, userID)
	if err != nil {
		return post, err
	}

	return post, nil
}

// GetUserVote retrieves the vote type (e.g., upvote, downvote) for a specific user, post, and comment.
func (db *Connection) GetUserVote(userID, postID, commentID int) (string, error) {
	var voteType string
	query := `SELECT vote_type FROM votes WHERE user_id = ? AND post_id = ? AND comment_id = ?`
	err := db.DB.QueryRow(query, userID, postID, commentID).Scan(&voteType)
	if err != nil {
		return "none", err
	}
	return voteType, nil
}

// SetVote sets or updates a user's vote (upvote/downvote) for a specific post or comment.
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

// GetPostVoteCounts retrieves the count of upvotes and downvotes for a specific post.
func (db *Connection) GetPostVoteCounts(postID int) (int, int) {
	var upvotes, downvotes int
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'upvote'`, postID).Scan(&upvotes)
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE post_id = ? AND vote_type = 'downvote'`, postID).Scan(&downvotes)
	return upvotes, downvotes
}

// GetCommentVoteCounts retrieves the count of upvotes and downvotes for a specific comment.
func (db *Connection) GetCommentVoteCounts(commentID int) (int, int) {
	var upvotes, downvotes int
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'upvote'`, commentID).Scan(&upvotes)
	db.DB.QueryRow(`SELECT COUNT(*) FROM votes WHERE comment_id = ? AND vote_type = 'downvote'`, commentID).Scan(&downvotes)
	return upvotes, downvotes
}

// DeletePost deletes a post from the database if the user is the author of the post.
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
