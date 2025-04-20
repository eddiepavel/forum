package auth

import (
	"encoding/json"
	"fmt"
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/render"
	"net/http"
	"strconv"
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

func PostVote(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
		voteType := r.FormValue("vote_type")

		// Register the vote in the database
		err := app.DB.SetVote(user.ID, postID, commentID, voteType)
		if err != nil {
			http.Error(w, "Failed to register vote", http.StatusInternalServerError)
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
