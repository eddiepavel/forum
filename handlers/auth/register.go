package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/helpers"
	"forum-app/helpers/flash"
	"forum-app/helpers/validator"
	"forum-app/render"
	"net/http"
)

func GetRegister(w http.ResponseWriter, r *http.Request) {
	view, err := render.PrepareView("register", r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = view.Render(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
func StoreRegister(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		inputs := map[string][]interface{}{
			"email":            {"required", "string", "email", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"username":         {"required", "string", app.DB.CheckUserExists(r.FormValue("email"), r.FormValue("username"))},
			"password":         {"required", "string"},
			"confirm_password": {"required", "string", "same:password"},
		}

		valid, errors := validator.ValidateRequest(r, inputs)

		if !valid {
			flash.HandleMessages(w, r, errors, r.Header.Get("Referer"), "error")
			return
		}

		userEmail := r.FormValue("email")
		userName := r.FormValue("username")
		userPassword, err := helpers.HashPassword(r.FormValue("password"))

		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = app.DB.RegisterUser(userEmail, userName, userPassword)
		if err != nil {
			app.Logger.Info("Failed to hash password", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}

}
