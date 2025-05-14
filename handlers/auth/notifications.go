package auth

import (
	"encoding/json"
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/render"
	"net/http"
)

func GetUnreadNotifs(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		data, err := app.DB.GetUnreadNotifications(user.ID)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}
		if data == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"notifications": []}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"notifications": data,
		})
	}
}

func MarkAsReadNotifs(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		err := app.DB.MarkAsReadNotifications(user.ID)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetNotifications(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		data, err := app.DB.GetNotifications(user.ID)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}
		if data == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"notifications": []}`))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"notifications": data,
		})
	}
}
