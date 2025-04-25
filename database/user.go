package database

import (
	"fmt"
	"forum-app/models"
	"time"
)

// CheckUserExists returns a function that checks if a user with the given email or username exists in the database.
func (db *Connection) CheckUserExists(email, username string) func(interface{}) error {
	return func(value interface{}) error {
		var emailExists, usernameExists bool
		query := `SELECT 
            EXISTS(SELECT 1 FROM user WHERE email = ? LIMIT 1) AS emailExists,
            EXISTS(SELECT 1 FROM user WHERE username = ? LIMIT 1) AS usernameExists;`

		err := db.DB.QueryRow(query, email, username).Scan(&emailExists, &usernameExists)
		if err != nil {
			return fmt.Errorf("database error: %v", err)
		}

		if emailExists && usernameExists {
			return fmt.Errorf("user with email %s and username %s already exist", email, username)
		} else if emailExists {
			return fmt.Errorf("user with email %s already exists", email)
		} else if usernameExists {
			return fmt.Errorf("user with username %s already exists", username)
		}

		return nil
	}
}

// RegisterUser inserts a new user into the database with the provided email, username, and hashed password.
func (db *Connection) RegisterUser(email, username, hashedPassword string) error {
	query := `INSERT INTO user (email, username, password, createdAt)
	          VALUES (?, ?, ?, ?)`

	_, err := db.DB.Exec(query, email, username, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// GetUserByEmail retrieves a user from the database by their email address.
func (db *Connection) GetUserByEmail(email string) (models.Users, error) {
	query := `SELECT * FROM user WHERE email = ? LIMIT 1;`
	var user models.Users

	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Is_Admin, &user.CreatedAt)

	return user, err
}

// GetUserById retrieves a user from the database by their ID, excluding the password from the result.
func (db *Connection) GetUserById(id int) (models.Users, error) {
	query := `SELECT * FROM user WHERE id = ? LIMIT 1;`
	var user models.Users

	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Is_Admin, &user.CreatedAt)
	user.Password = ""
	return user, err
}
