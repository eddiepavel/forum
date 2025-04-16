package render

import (
	"fmt"
	"forum-app/app"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"html/template"
	"net/http"
	"strconv"
)

var files = []string{
	"./assets/base.html",
	"./assets/partials/nav.html",
	"./assets/partials/create.html",
	"./assets/partials/home.html",
	"./assets/partials/posts.html",
	"./assets/partials/view.html",
	"./assets/partials/wip.html",
	"./assets/partials/login.html",
	"./assets/partials/register.html"}

var categories = []string{
	"General",
	"Technology",
	"Entertainment",
	"Sports",
	"News",
	"Anouncements",
	"Other"}

func (view *View) Render(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(view.Path...)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, view.Data)
	if err != nil {
		return err
	}
	return nil
}

func PrepareView(source string, r *http.Request, app *app.Application) (View, error) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	session := r.Context().Value("user_session").(*session.Session)
	redirect := r.URL.Query().Get("redirect")

	data := models.PageData{}
	if ok && user != nil {
		data = models.PageData{User: user, Session: session}
	} else {
		data = models.PageData{User: nil, Session: session}
	}

	data.Data = make(map[string]interface{})

	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}

	// Retrieve flash messages
	if flash, exists := session.GetFlash("error"); exists {
		data.Data["error"] = flash
	}

	if source == "home" {
		posts, err := app.DB.GetPostsForHome(1, r.URL.Query().Get("category"), user)
		if err != nil {
			return View{}, err
		}
		data.Data["posts"] = posts
	}

	if source == "view" {
		postID := r.URL.Query().Get("id")
		id, err := strconv.Atoi(postID)
		if err != nil {
			return View{}, fmt.Errorf("invalid post ID: %v", err)
		}
		if postID == "" {
			return View{}, fmt.Errorf("post ID is required")
		}
		post, err := app.DB.GetPostByID(id)
		if err != nil {
			return View{}, err
		}
		data.Data["post"] = post
	}

	data.Source = source

	if source == "create" || source == "home" {
		data.Data["categories"] = categories
	}

	view := View{
		Name: source,
		Data: &data,
	}

	if redirect != "" {
		data.Redirect = redirect
	}

	view.Path = files

	return view, nil
}
