package render

import "forum-app/models"

type View struct {
	Name string
	Path []string
	Data *models.PageData
}
