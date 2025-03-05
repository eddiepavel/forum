package app

import (
	"forum-app/database"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	DB     *database.Connection
	Logger *slog.Logger
}
