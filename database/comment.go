package database

import (
	"forum-app/utils"
	"time"
)

func (db *Connection) SetComment(postID int, content, author string) error {
	// Sanitize comment
	cleanContent, err := utils.SanitizeComment(content)
	if err != nil {
		return err
	}

	query := `INSERT INTO comment(content, author, post_id, time)
				VALUES(?, ?, ?, ?)`

	_, err = db.DB.Exec(query, cleanContent, author, postID, time.Now().Format("2006-01-02 15:04:05"))
	return err
}
