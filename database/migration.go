package database

import (
	"fmt"
	"os"
)

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
