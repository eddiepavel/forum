package auth_user

import (
	"forum-app/middleware"
	"forum-app/models"
	"net/http"
)

// AuthUser retrieves the authenticated user from the request context.
func AuthUser(r *http.Request) *models.Users {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	if ok && user != nil {
		return user
	}
	return nil
}
