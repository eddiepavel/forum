package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/handlers/forum"
	"forum-app/session"
	"html/template"
	"net/http"
)

// GetCreate is a handler function that returns the create forum page.
func GetCreate(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("user_session").(*session.Session)
	data := forum.PageData{Session: session}

	t, err := template.ParseFiles("./assets/create.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
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
