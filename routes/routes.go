package routes

import (
	"forum-app/app"
	"forum-app/handlers/auth"
	"forum-app/handlers/forum"
	"forum-app/middleware"
	"net/http"
)

func Web(app *app.Application) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/", middleware.ChainMiddleware(forum.GetRedirect(app), []string{"auth"}, app))
	mux.HandleFunc("GET /home", middleware.ChainMiddleware(forum.GetHome(app), []string{"auth"}, app))
	mux.HandleFunc("GET /login", middleware.ChainMiddleware(auth.GetLogin(app), []string{}, app))
	mux.HandleFunc("POST /login", middleware.ChainMiddleware(auth.PostLogin(app), []string{}, app))
	mux.HandleFunc("GET /register", middleware.ChainMiddleware(auth.GetRegister(app), []string{}, app))
	mux.HandleFunc("POST /register", middleware.ChainMiddleware(auth.StoreRegister(app), []string{}, app))
	mux.HandleFunc("GET /logout", middleware.ChainMiddleware(auth.Logout(app), []string{"auth"}, app))
	mux.HandleFunc("GET /create", middleware.ChainMiddleware(auth.GetCreate(app), []string{"auth"}, app))
	mux.HandleFunc("POST /create", middleware.ChainMiddleware(auth.PostCreate(app), []string{"auth"}, app))
	mux.HandleFunc("GET /view", middleware.ChainMiddleware(auth.GetView(app), []string{"auth"}, app))
	mux.HandleFunc("POST /view", middleware.ChainMiddleware(auth.PostView(app), []string{"auth"}, app))
	mux.HandleFunc("/wip", middleware.ChainMiddleware(forum.GetWIP(app), []string{"auth"}, app))

	return mux
}
