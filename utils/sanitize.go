package utils

import (
	"errors"
	"html"
	"strings"
)

func SanitizePost(title, category, content string) (string, string, string, error) {
	// Trim spaces
	title = strings.TrimSpace(title)
	category = strings.TrimSpace(category)
	content = strings.TrimSpace(content)

	// Escape HTML special characters
	title = html.EscapeString(title)
	category = html.EscapeString(category)
	content = html.EscapeString(content)

	// Validate lengths
	if len(title) < 1 || len(title) > 100 {
		return "", "", "", errors.New("title must be between 1 and 100 characters")
	}
	if len(category) < 1 || len(category) > 50 {
		return "", "", "", errors.New("category must be between 1 and 50 characters")
	}
	if len(content) < 1 || len(content) > 5000 {
		return "", "", "", errors.New("content must be between 1 and 5000 characters")
	}

	// Remove multiple spaces
	title = strings.Join(strings.Fields(title), " ")
	category = strings.Join(strings.Fields(category), " ")

	return title, category, content, nil
}

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
