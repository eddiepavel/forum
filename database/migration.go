package database

import (
	"fmt"
	"os"
)

// RunMigrations runs the SQL migration files in the migrations directory.
// It reads each file, executes the SQL commands, and prints the name of each file executed.
// If an error occurs during the process, it returns the error.
// It is assumed that the database connection is already established and available in the Connection struct.
func (db *Connection) RunMigrations() error {

	sqlFile, err := os.ReadDir("./migrations")

	if err != nil {
		fmt.Println("Failed to run migrations reason: ", err)
		os.Exit(1)
	}

	for _, file := range sqlFile {
		fileContents, _ := os.ReadFile("./migrations/" + file.Name())
		_, err := db.DB.Exec(string(fileContents))
		if err != nil {
			return fmt.Errorf("failed to run migrations reason: %w", err)
		}

		fmt.Println(file.Name())
	}

	return nil
}
