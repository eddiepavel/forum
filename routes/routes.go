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

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./assets/app"))))
	mux.HandleFunc("GET /{$}", middleware.ChainMiddleware(forum.GetHome(app), []string{"auth"}, app))
	mux.HandleFunc("GET /login", middleware.ChainMiddleware(auth.GetLogin, []string{}, app))
	mux.HandleFunc("POST /login", middleware.ChainMiddleware(auth.PostLogin(app), []string{}, app))
	mux.HandleFunc("GET /register", middleware.ChainMiddleware(auth.GetRegister, []string{}, app))
	mux.HandleFunc("POST /register", middleware.ChainMiddleware(auth.StoreRegister(app), []string{}, app))
	mux.HandleFunc("GET /logout", middleware.ChainMiddleware(auth.Logout(app), []string{"auth"}, app))

	return mux
}
