package render

import (
	"fmt"
	"forum-app/middleware"
	"forum-app/models"
	"forum-app/session"
	"html/template"
	"net/http"
	"strings"
)

var files = []string{
	"./assets/base.html",
	"./assets/partials/nav.html",
	"./assets/partials/create.html",
	"./assets/partials/home.html",
	"./assets/partials/posts.html",
	"./assets/partials/view.html",
	"./assets/partials/wip.html",
	"./assets/partials/error.html",
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

func PrepareView(source string, r *http.Request) (View, error) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.Users)
	session := r.Context().Value("user_session").(*session.Session)
	redirect := r.URL.Query().Get("redirect")

	data := models.PageData{}
	if ok && user != nil {
		data = models.PageData{User: user, Session: session}
	} else {
		data = models.PageData{User: nil, Session: session}
	}

	if strings.HasPrefix(source, "error:") {
		fmt.Println("Error page requested:", source)
		data.Data = make(map[string]interface{})
		if source == "error:404" {
			data.Data["message"] = "Page Not Found"
		}
		if source == "error:500" {
			data.Data["message"] = "Internal Server Error"
		}
		if source == "error:403" {
			data.Data["message"] = "Forbidden"
		}
		view := View{
			Name: source[:6],
			Data: &data,
		}
		data.Source = "error"
		view.Path = files
		return view, nil
	}

	posts, err := models.GetData(source, r)
	if err != nil {
		return View{}, err
	}

	data.Data = make(map[string]interface{})
	if source == "home" {
		data.Data["posts"] = posts
	} else if source == "view" {
		data.Data["post"] = posts
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
