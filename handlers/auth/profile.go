package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/render"
	"net/http"
	"strconv"
)

func GetUserProfile(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("profile", r, app)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}
		err = view.Render(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func RedirectToProf(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
		if !ok || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		id := strconv.Itoa(user.ID)
		http.Redirect(w, r, "/profile/"+id, http.StatusSeeOther)
	}
}
