package database

import "forum-app/models"

func (db *Connection) CheckVoteNotification(targetID, actorID, postID int, voteType string) (int, error) {
	return -1, nil
}

func (db *Connection) DeleteNotification(id int) error {
	return nil
}

func (db *Connection) GetUnreadNotifications(userID int) ([]models.Notification, error) {
	return nil, nil
}

func (db *Connection) MarkAsReadNotifications(userID int) error {
	return nil
}

func (db *Connection) GetNotifications(userID int) ([]models.Notification, error) {
	return nil, nil
}

func (db *Connection) SetNotification(targetID int, actorID int, postID int, voteType string, comment string) error {
	return nil
}
