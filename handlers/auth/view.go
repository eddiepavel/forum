package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

func GetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("view", r)
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
}

func PostView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		for key, values := range r.Form {
			fmt.Printf("Key: %s, Values: %v\n", key, values)
		}

		redirect := r.FormValue("redirect")

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
