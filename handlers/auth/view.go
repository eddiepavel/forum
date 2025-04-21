package auth

import (
	"encoding/json"
	"fmt"
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/render"
	"html/template"
	"net/http"
	"strconv"
)

func GetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := render.PrepareView("view", r, app)
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

func PostView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		comment := r.FormValue("comment")
		postId := r.FormValue("post_id")
		authorId := r.FormValue("author_id")

		err := app.DB.SetComment(postId, comment, authorId)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		redirect := r.FormValue("redirect")

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}

func DeletePost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || postID <= 0 {
			render.RenderError(w, r, fmt.Errorf("invalid post ID"))
			return
		}

		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		err = app.DB.DeletePost(postID, user.ID)
		if err != nil {
			tmpl, err := template.ParseFiles("./assets/error.html")
			if err != nil {
				render.RenderError(w, r, err)
				return
			}
			data := err.Error()
			err = tmpl.Execute(w, data)
			if err != nil {
				render.RenderError(w, r, err)
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func PostVote(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			render.RenderError(w, r, err)
			return
		}

		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
		voteType := r.FormValue("vote_type")

		// Register the vote in the database
		err := app.DB.SetVote(user.ID, postID, commentID, voteType)
		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		// Fetch the updated vote counts and user's vote state
		var upvotes, downvotes int
		var userVote string
		if postID != 0 {
			upvotes, downvotes = app.DB.GetPostVoteCounts(postID)
			userVote = app.DB.GetUserVote(user.ID, postID, 0)
		} else if commentID != 0 {
			upvotes, downvotes = app.DB.GetCommentVoteCounts(commentID)
			userVote = app.DB.GetUserVote(user.ID, 0, commentID)
		}

		// Return the updated data as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"upvotes":   upvotes,
			"downvotes": downvotes,
			"user_vote": userVote,
		})
	}
}
