package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

// GetCreate is a handler function that returns the create forum page.
func GetCreate(w http.ResponseWriter, r *http.Request) {
	view, err := render.PrepareView("create", r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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

func PostCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Print form values to the terminal
		for key, values := range r.Form {
			fmt.Printf("Key: %s, Values: %v\n", key, values)
		}

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
