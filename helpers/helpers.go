package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// CompareHashAndPassword compares a hashed password with a plain-text one.
func CompareHashAndPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// GenerateToken generates a random 128-bit token encoded as a hexadecimal string.
func GenerateToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// DdSessionTimeSeconds calculates the remaining session time in seconds based on a target date.
func DdSessionTimeSeconds(date string) int {

	layout := "2006-01-02 15:04:05"

	targetTime, err := time.Parse(layout, date)
	if err != nil {
		return -1
	}

	currentTime := time.Now()

	maxAge := int(targetTime.Sub(currentTime).Seconds())

	if maxAge < 0 {
		maxAge = 0
	}

	return maxAge
}

// CompareDatesLess checks if the first date is earlier than the second date.
func CompareDatesLess(date1 time.Time, date2 string) bool {
	layout := "2006-01-02 15:04:05"

	time2, err := time.Parse(layout, date2)
	if err != nil {
		return false
	}

	return date1.Before(time2)
}

// SanitizePost sanitizes and validates the title and content of a post.
func SanitizePost(title, content string) (string, string, error) {
	// Trim spaces
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	// Escape HTML special characters
	title = html.EscapeString(title)
	content = html.EscapeString(content)

	// Validate lengths
	if len(title) < 1 || len(title) > 200 {
		return "", "", errors.New("title must be between 1 and 200 characters")
	}
	if len(content) < 1 || len(content) > 5000 {
		return "", "", errors.New("content must be between 1 and 5000 characters")
	}

	// Remove multiple spaces
	title = strings.Join(strings.Fields(title), " ")

	return title, content, nil
}

// SanitizeComment sanitizes and validates the content of a comment.
func SanitizeComment(content string) (string, error) {
	// Trim spaces
	content = strings.TrimSpace(content)

	// Escape HTML special characters
	content = html.EscapeString(content)

	// Validate length
	if len(content) < 1 || len(content) > 1000 {
		return "", errors.New("comment must be between 1 and 1000 characters")
	}

	return content, nil
}

// Beautify converts error messages into user-friendly messages.
func Beautify(err error) string {
	if err == nil {
		return ""
	}
	if strings.Contains(err.Error(), "duplicate") {
		return "This record already exists. Please use a unique value."
	}
	if strings.Contains(err.Error(), "not found") {
		return "The requested item could not be found. Please check your input."
	}
	if strings.Contains(err.Error(), "invalid") {
		return "The input provided is invalid. Please correct it and try again."
	}
	if strings.Contains(err.Error(), "redirect URL") {
		return "The redirect URL is invalid. Please contact support if the issue persists."
	}
	if strings.Contains(err.Error(), "validation rule") {
		return "There was an issue with the validation rules. Please contact support."
	}
	if strings.Contains(err.Error(), "length") {
		return "The input length is invalid. Please adhere to the specified limits."
	}
	if strings.Contains(err.Error(), "sql: no rows") {
		return "No records found matching your criteria. Please check your input."
	}
	return "An unexpected error occurred: " + err.Error()
}

// BeautifyMessage formats a message by splitting on underscores and capitalizing the first word.
func BeautifyMessage(message string) string {
	newMessage := strings.Split(message, "_")

	if len(newMessage) > 1 {
		newMessage[0] = strings.Title(newMessage[0])
		return strings.Join(newMessage, " ")
	} else {
		newMessage = strings.Split(message, " ")
		return strings.Title(newMessage[0]) + " " + strings.Join(newMessage[1:], " ")
	}
}
