package middleware

import (
	"fmt"
	"forum-app/app"
	"net/http"
	"strconv"
)

func EventMiddleware(next http.HandlerFunc, app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		////Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max memory (adjust as needed)
		if err != nil {
			fmt.Println("ParseMultipartForm error:", err)
		}

		// Extract values
		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		voteType := r.FormValue("vote_type")
		comment := r.FormValue("comment")
		targetID, _ := strconv.Atoi(r.FormValue("post_author_id"))
		actorID, _ := strconv.Atoi(r.FormValue("action_author_id"))

		if actorID == targetID {
			next(w, r)
			return
		}

		if voteType != "" {
			status, err := app.DB.CheckVoteNotification(targetID, actorID, postID, voteType)
			if err != nil {
				fmt.Println(err)
			}
			if status != -1 {
				err := app.DB.DeleteNotification(status)
				if err != nil {
					fmt.Println(err)
				}
				next(w, r)
				return
			}
		}
		err = app.DB.SetNotification(targetID, actorID, postID, voteType, comment)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("Post ID: ", postID)
		//fmt.Println("Vote Type: ", voteType)
		//fmt.Println("Comment: ", comment)
		//fmt.Println("Target ID: ", targetID)
		//fmt.Println("Actor ID: ", actorID)

		//fmt.Println("------ Incoming Request ------")
		//// Print request body if present
		//fmt.Println("\nBody:")
		//if r.Body != nil {
		//	// Read and restore body so it remains usable
		//	bodyBytes, _ := io.ReadAll(r.Body)
		//	fmt.Println(string(bodyBytes))
		//	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		//}
		//
		//fmt.Println("\n-----------------------------")

		next(w, r)
	}
}
