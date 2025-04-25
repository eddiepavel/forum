package validator

import (
	"errors"
	"fmt"
	"forum-app/helpers"
	"net/mail"
	"regexp"
	"strconv"
)

// Exists checks if a value exists in the specified table and column in the database.
func (v *Validator) Exists(value interface{}, table, column string) error {
	if v.app == nil || v.app.DB == nil || v.app.DB.DB == nil {
		return errors.New("database connection is not available")
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
	var count int
	err := v.app.DB.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return errors.New("database error: " + err.Error())
	}
	if count == 0 {
		return errors.New(fmt.Sprintf("value '%v' does not exist in %s.%s", value, table, column))
	}
	return nil
}

// ValidateLoginAttempt validates a user's login attempt by checking email and password.
func (v *Validator) ValidateLoginAttempt(email, password string) error {
	user, err := v.app.DB.GetUserByEmail(email)
	if err != nil || helpers.CompareHashAndPassword(user.Password, password) != nil {
		return errors.New("invalid email or password")
	}
	return nil
}

// ValidateString checks if the value is a valid string.
func (v *Validator) ValidateString(value interface{}, key string) error {
	_, ok := value.(string)
	if !ok {
		return errors.New(key + " value is not a valid string")
	}
	return nil
}

// ValidateInt checks if the value is a valid integer.
func (v *Validator) ValidateInt(value interface{}, key string) error {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return nil
	case string:
		if _, err := strconv.Atoi(v); err == nil {
			return nil
		}
	}
	return errors.New(key + " value is not a valid integer")
}

// ValidateEmail checks if the value is a valid email address.
func (v *Validator) ValidateEmail(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}
	_, err := mail.ParseAddress(str)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

// Required checks if the value is not empty.
func (v *Validator) Required(value interface{}, key string) error {
	if value == "" {
		return errors.New(key + " is required")
	}
	return nil
}

// ValidatePassword validates a password based on modular rules.
// It ensures the password has at least 8 characters and at least one number.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check if the password contains at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	if !hasNumber(password) {
		return errors.New("password must contain at least one number")
	}

	hasAlpha := regexp.MustCompile(`[a-zA-Z]`).MatchString
	if !hasAlpha(password) {
		return errors.New("password must contain at least one letter")
	}

	// Add more rules here if needed
	return nil
}
