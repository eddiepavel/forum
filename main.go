package main

import (
	"flag"
	"forum-app/app"
	"forum-app/database"
	"forum-app/ratelimiter"
	"forum-app/routes"
	"forum-app/session"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "app.db", "Database file name sqlite3")

	flag.Parse()

	db, err := database.NewConnection(*dbName)

	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	defer db.DB.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	session := session.NewSessionStore(1*time.Hour, 1*time.Hour)

	rl := ratelimiter.NewRateLimiter(100, 1*time.Minute)

	app := &app.Application{
		DB:          db,
		Logger:      logger,
		Session:     session,
		RateLimiter: rl,
	}

	server := http.Server{
		Addr:    *addr,
		Handler: routes.Web(app),
	}

	logger.Info("starting server", "addr", *addr)

	server.ListenAndServe()

	os.Exit(1)
}
