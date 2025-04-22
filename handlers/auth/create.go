package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/helpers/validator"
	"forum-app/render"
	"net/http"
)

// GetCreate returns an HTTP handler function for rendering the create post page.
func GetCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("create", r, app)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		if view.Data.User == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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

// PostCreate handles the creation of a new forum post by validating input and saving it to the database.
func PostCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		// Define validation rules
		inputs := map[string][]interface{}{
			"title":       {"required", "string"},
			"description": {"required", "string"},
			"categories":  {"sometimes", "string"},
			"user_id":     {"required", "exists:user,id", "string"},
		}

		// Validate the request
		valid, errors := validator.ValidateRequest(r, inputs, app)
		if !valid {
			cookie, err := r.Cookie("session")
			if err != nil {
				render.RenderError(w, r, err)
				return
			}
			session, _ := app.Session.GetSession(cookie.Value)
			session.SetFlash("error", errors)
			http.Redirect(w, r, "/create", http.StatusFound)
			return
		}

		// Extract validated form values
		title := r.FormValue("title")
		content := r.FormValue("description")
		categories := r.FormValue("categories")
		author := r.FormValue("user_id")

		if categories == "" {
			categories = "General"
		}

		// Save the post to the database
		err = app.DB.SetPost(title, content, author, categories)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		// Redirect to the home page
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
