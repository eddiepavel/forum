package forum

import (
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

func GetHome(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view := render.PrepareView("home", r)

		err := view.Render(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
