package app

import (
	"forum-app/database"
	"forum-app/session"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	DB     *database.Connection
	Logger *slog.Logger
	Session *session.SessionStore
}
