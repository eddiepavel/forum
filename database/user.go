package database

import (
	"fmt"
	"forum-app/models"
	"time"
)

func (db *Connection) CheckUserExists(email, username string) func(interface{}) error {
	return func(value interface{}) error {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM user WHERE email = ? OR username = ?  LIMIT 1);`

		err := db.DB.QueryRow(query, email, username).Scan(&exists)
		if err != nil {
			return fmt.Errorf("database error: %v", err)
		}

		if exists {
			return fmt.Errorf("user with email %s already exists", email)
		}

		return nil
	}
}

func (db *Connection) RegisterUser(email, username, hashedPassword string) error {
	query := `INSERT INTO user (email, username, password, createdAt)
	          VALUES (?, ?, ?, ?)`

	_, err := db.DB.Exec(query, email, username, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func (db *Connection) GetUserByEmail(email string) (models.Users, error) {
	query := `SELECT * FROM user WHERE email = ? LIMIT 1;`
	var user models.Users

	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Is_Admin, &user.CreatedAt)

	return user, err
}

func (db *Connection) GetUserById(id int) (models.Users, error) {
	query := `SELECT * FROM user WHERE id = ? LIMIT 1;`
	var user models.Users

	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Is_Admin, &user.CreatedAt)

	return user, err
}
