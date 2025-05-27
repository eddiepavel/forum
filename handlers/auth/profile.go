package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"net/http"
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
