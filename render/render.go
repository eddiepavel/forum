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
	user, session := getUserAndSession(r)
	data := initializePageData(user, session)

	handleFlashMessages(session, &data)

	switch source {
	case "home":
		if err := handleHomePage(r, app, user, &data); err != nil {
			return View{}, err
		}
	case "view":
		if err := handleViewPage(r, app, &data); err != nil {
			return View{}, err
		}
	}

	if source == "create" || source == "home" {
		setCategories(&data)
	}

	data.Source = source

	view := View{
		Name: source,
		Data: &data,
	}

	redirect := r.URL.Query().Get("redirect")
	if redirect != "" {
		data.Redirect = redirect
	}

	view.Path = files

	return view, nil
}

func getUserAndSession(r *http.Request) (*models.Users, *session.Session) {
	user, _ := r.Context().Value(middleware.UserKey).(*models.Users)
	session := r.Context().Value("user_session").(*session.Session)
	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}
	return user, session
}

func initializePageData(user *models.Users, session *session.Session) models.PageData {
	if user != nil {
		return models.PageData{User: user, Session: session, Data: make(map[string]interface{})}
	}
	return models.PageData{User: nil, Session: session, Data: make(map[string]interface{})}
}

func handleFlashMessages(session *session.Session, data *models.PageData) {
	if flash, exists := session.GetFlash("error"); exists {
		data.Data["error"] = flash
	}
}

func handleHomePage(r *http.Request, app *app.Application, user *models.Users, data *models.PageData) error {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return fmt.Errorf("invalid page number: %v", err)
	}
	filter := r.URL.Query().Get("category")
	fmt.Println("Filter 1:", filter)
	totalPosts, err := app.DB.GetTotalPostCount(r.URL.Query().Get("category"), user)
	fmt.Println("Total Posts:", totalPosts)
	if err != nil {
		return err
	}

	const pageSize = 10
	totalPages := (totalPosts + pageSize - 1) / pageSize
	if (pageNum > totalPages || pageNum < 1) && totalPosts != 0 {
		return fmt.Errorf("Couldn't find page %d", pageNum)
	}

	if totalPosts == 0 {
		data.Data["posts"] = nil
	} else {
		posts, err := app.DB.GetPostsForHome(pageNum, r.URL.Query().Get("category"), user)
		if err != nil {
			return err
		}
		data.Data["posts"] = posts
		data.Data["totalPosts"] = totalPosts
		data.Data["fromPosts"] = 1 + ((pageNum - 1) * pageSize)
		data.Data["toPosts"] = len(posts) + ((pageNum - 1) * pageSize)
	}
	return nil
}

func handleViewPage(r *http.Request, app *app.Application, data *models.PageData) error {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		return fmt.Errorf("post ID is required")
	}
	id, err := strconv.Atoi(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %v", err)
	}
	userID := -1
	if data.User != nil {
		userID = data.User.ID
	}
	post, err := app.DB.GetPostByID(id, userID)
	if err != nil {
		return err
	}
	data.Data["post"] = post
	return nil
}

func setCategories(data *models.PageData) {
	data.Data["categories"] = categories
}
