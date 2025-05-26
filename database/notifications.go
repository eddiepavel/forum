package database

import (
	"database/sql"
	"forum-app/models"
	"time"
)

func (db *Connection) CheckVoteNotification(targetID, actorID, postID int, voteType string) (int, error) {
	var id int
	query := `SELECT id FROM notification WHERE target_id = ? AND actor_id = ? AND post_id = ? AND type = ?`
	err := db.DB.QueryRow(query, targetID, actorID, postID, voteType).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return id, nil

}

func (db *Connection) DeleteNotification(id int) error {
	query := `DELETE FROM notification WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}

func (db *Connection) GetUnreadNotifications(userID int) ([]models.Notification, error) {
	query := `SELECT n.id, n.target_id, n.actor_id, n.post_id, n.type, n.content, n.read, n.time FROM notification n WHERE n.target_id = ? AND n.read = 0 ORDER BY time DESC`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var target, actor int
		var timeRaw time.Time
		err := rows.Scan(
			&notification.ID,
			&target,
			&actor,
			&notification.PostID,
			&notification.Type,
			&notification.Content,
			&notification.Read,
			&timeRaw,
		)
		if err != nil {
			return nil, err
		}
		notification.Target, err = db.GetUserById(target)
		notification.Target.Password = ""
		if err != nil {
			return nil, err
		}
		notification.Actor, err = db.GetUserById(actor)
		notification.Actor.Password = ""
		if err != nil {
			return nil, err
		}
		notification.Time = timeRaw.Format("2006-01-02 15:04:05")
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (db *Connection) MarkAsReadNotifications(userID int) error {
	query := `UPDATE notification SET read = 1 WHERE target_id = ? AND read = 0`
	_, err := db.DB.Exec(query, userID)
	return err
}

func (db *Connection) GetNotifications(userID int) ([]models.Notification, error) {
	query := `SELECT n.id, n.target_id, n.actor_id, n.post_id, n.type, n.content, n.read, n.time FROM notification n WHERE target_id = ? ORDER BY time DESC`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var target, actor int
		var timeRaw time.Time
		err := rows.Scan(
			&notification.ID,
			&target,
			&actor,
			&notification.PostID,
			&notification.Type,
			&notification.Content,
			&notification.Read,
			&timeRaw,
		)
		if err != nil {
			return nil, err
		}
		notification.Target, err = db.GetUserById(target)
		if err != nil {
			return nil, err
		}
		notification.Actor, err = db.GetUserById(actor)
		if err != nil {
			return nil, err
		}
		notification.Time = timeRaw.Format("2006-01-02 15:04:05")
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (db *Connection) SetNotification(targetID int, actorID int, postID int, voteType string, comment string) error {
	var notifType, content string
	if comment != "" {
		notifType = "comment"
		if len(comment) > 20 {
			content = comment[:20] + "..."
		} else {
			content = comment
		}
	} else {
		notifType = voteType
		content = ""
	}
	query := `INSERT INTO notification(target_id, actor_id, post_id, type, content, time)
				VALUES(? ,? ,? ,? ,? ,?)`
	_, err := db.DB.Exec(query, targetID, actorID, postID, notifType, content, time.Now().Format("2006-01-02 15:04:05"))
	return err
}
