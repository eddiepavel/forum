package auth_user

import (
	"forum-app/middleware"
	"forum-app/models"
	"net/http"
)

func AuthUser(r *http.Request) *models.Users {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	if ok && user != nil {
		return user
	}
	return nil
}
