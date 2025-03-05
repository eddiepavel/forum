package forum

import "forum-app/models"

type PageData struct {
	Data []interface{}
	User *models.Users
}
