package forum

import (
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

// GetWIP returns an HTTP handler function for rendering the "Work In Progress" page.
func GetWIP(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("wip", r, app)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		err = view.Render(w, r)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
