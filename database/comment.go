package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-app/helpers"
	"time"
)

// Add a comment to a post, in the db
func (db *Connection) SetComment(postID, content, author string) error {
	// Sanitize comment
	cleanContent, err := helpers.SanitizeComment(content)
	if err != nil {
		return err
	}

	query := `INSERT INTO comment(content, author, post_id, time)
				VALUES(?, ?, ?, ?)`

	_, err = db.DB.Exec(query, cleanContent, author, postID, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func (db *Connection) UpdateComment(id int, content string, userID int) error {
	var authorID int
	err := db.DB.QueryRow(`SELECT author FROM comment WHERE id = ?`, id).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("comment not found")
		}
	}
	if authorID != userID {
		return errors.New("403 Forbidden: You are not the author of this comment")
	}

	cleanContent, err := helpers.SanitizeComment(content)
	if err != nil {
		return err
	}

	query := `UPDATE comment SET content = ? WHERE id = ?`
	result, err := db.DB.Exec(query, cleanContent, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (db *Connection) DeleteComment(commentID, userID int) error {
	var authorID int
	err := db.DB.QueryRow(`SELECT author FROM comment WHERE id = ?`, commentID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("comment not found")
		}
		return err
	}

	if authorID != userID {
		return errors.New("403 Forbidden: You are not the author of this comment")
	}

	// Proceed with deletion
	query := `DELETE FROM comment WHERE id = ? AND author = ?`
	result, err := db.DB.Exec(query, commentID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}
