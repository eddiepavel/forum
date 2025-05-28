package auth

import (
	"fmt"
	"forum-app/app"
	"forum-app/helpers/validator"
	"forum-app/render"
	"io"
	"net/http"
	"os"
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
		// Define validation rules
		inputs := map[string][]interface{}{
			"title":       {"required", "string"},
			"description": {"required", "string"},
			"categories":  {"sometimes", "string"},
			"user_id":     {"required", "exists:user,id", "string"},
		}

		var imagePath string
		// Handle image upload
		if r.Method == "POST" {
			file, handler, err := r.FormFile("image")
			var imagePath string
			if err == nil {
				defer file.Close()

				// Validate file type
				validTypes := map[string]bool{
					"image/jpeg": true,
					"image/png":  true,
					"image/gif":  true,
					"image/jpg":  true,
					"image/webp": true,
				}

				// Validate file size
				const maxFileSize = 20 * 1024 * 1024 // 20 MB
				if handler.Size > maxFileSize {
					http.Error(w, "Image is too large. Max allowed size is 20 MB.", http.StatusBadRequest)
					return
				}

				// Get the content type of the file
				buffer := make([]byte, 512)
				_, err = file.Read(buffer)
				if err != nil {
					render.RenderError(w, r, err)
					return
				}

				contentType := http.DetectContentType(buffer)
				if !validTypes[contentType] {
					http.Error(w, "Invalid image type. Allowed types are JPEG, PNG, and GIF.", http.StatusBadRequest)
					return
				}

				// Reset file pointer for saving
				_, err = file.Seek(0, 0)
				if err != nil {
					render.RenderError(w, r, err)
					return
				}

				// Save the image to ./uploads folder
				imagePath = "./uploads/" + handler.Filename
				dst, err := os.Create(imagePath)
				if err != nil {
					render.RenderError(w, r, err)
					return
				}
				defer dst.Close()

				_, err = io.Copy(dst, file)
				if err != nil {
					render.RenderError(w, r, err)
					return
				}
			}
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
		postID := r.FormValue("post_id")

		if categories == "" {
			categories = "General"
		}
		var err error
		// Save the post to the database
		if r.Method == "POST" {
			err = app.DB.SetPost(title, content, author, categories, imagePath)
		} else {
			err = app.DB.UpdatePost(title, content, author, categories, postID)
		}

		if err != nil {
			render.RenderError(w, r, err)
			return
		}

		if r.Method == "POST" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/view?id="+postID, http.StatusSeeOther)
	}
}
