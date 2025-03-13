package forum

import (
	"forum-app/models"
	"forum-app/session"
)

type PageData struct {
	Data []interface{}
	User *models.Users
	Session *session.Session
}
