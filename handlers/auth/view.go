package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/render"
	"net/http"
)

func GetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("view", r, app)
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

// Key: comment, Values: [this is a test comment]
// Key: post_id, Values: [1]
// Key: csrf, Values: [d08468e6abe379fbdb7262aa2c4f3a9f]
// Key: author_id, Values: [1]

func PostView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		comment := r.FormValue("comment")
		postId := r.FormValue("post_id")
		authorId := r.FormValue("author_id")

		err := app.DB.SetComment(postId, comment, authorId)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		redirect := r.FormValue("redirect")

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
