package forum

import (
	"errors"
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

// GetHome returns an HTTP handler function for rendering the home page.
func GetHome(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("home", r, app)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		err = view.Render(w, r)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

// GetRedirect returns an HTTP handler function for handling redirects to specific paths.
func GetRedirect(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		render.RenderError(w, r, errors.New("Page not found"))
	}
}
