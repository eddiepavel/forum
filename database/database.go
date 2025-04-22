package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	DB *sql.DB
}

// Opens a connection to the database and runs migrations if the database file does not exist.
// If the database file already exists, it simply opens the connection.
// Returns a pointer to the Connection struct and an error if any occurs.
// The database file is located at "././database/file/" + name.
func NewConnection(name string) (*Connection, error) {

	exists, _ := os.Stat("././database/file/" + name)

	db, err := sql.Open("sqlite3", "././database/file/"+name)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if exists == nil {
		fmt.Println("Running migrations...")
		connection := Connection{DB: db}
		err := connection.RunMigrations()

		if err != nil {
			if _, err := os.Stat("././database/file/" + name); errors.Is(err, os.ErrExist) {
				os.Remove("././database/file/" + name)
			}
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return &Connection{DB: db}, nil
}
