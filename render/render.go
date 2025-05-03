package render

import (
	"errors"
	"forum-app/app"
	"forum-app/helpers"
	"html/template"
	"net/http"
)

// Render renders the view using the provided HTTP response writer and request.
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

// RenderError renders an error page with the provided error message.
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	tmpl, parseErr := template.ParseFiles("./assets/error.html")
	if parseErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	data := helpers.Beautify(err)
	execErr := tmpl.Execute(w, data)
	if execErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// PrepareView prepares the view data and structure based on the source and request.
func PrepareView(source string, r *http.Request, app *app.Application) (View, error) {
	user, session := getUserAndSession(r)
	data := initializePageData(user, session)

	// Check for errors in the context
	if flash, exists := session.GetFlash("csrf_error"); exists {
		return View{}, errors.New(flash.(string))
	}

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
