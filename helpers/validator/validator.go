package validator

import (
	"errors"
	"fmt"
	"forum-app/app"
	"forum-app/helpers"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

type Validator struct {
	app *app.Application
}

func NewValidator(app *app.Application) *Validator {
	return &Validator{app: app}
}

func (v *Validator) ValidateString(value interface{}, key string) error {
	_, ok := value.(string)
	if !ok {
		return errors.New(key + " value is not a valid string")
	}
	return nil
}

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

func (v *Validator) Required(value interface{}, key string) error {
	if value == "" {
		return errors.New(key + " is required")
	}
	return nil
}

func (v *Validator) ValidateInput(value interface{}, rules []interface{}, key string, hold map[string]interface{}) error {
	for _, rule := range rules {
		switch rule := rule.(type) {
		case string: // Standard validation rules
			switch {
			case rule == "string":
				if err := v.ValidateString(value, key); err != nil {
					return err
				}
			case rule == "int":
				if err := v.ValidateInt(value, key); err != nil {
					return err
				}
			case rule == "email":
				if err := v.ValidateEmail(value); err != nil {
					return err
				}
			case rule == "required":
				if err := v.Required(value, key); err != nil {
					return err
				}
			case rule == "sometimes":
				// Skip validation if the field is not present
				if value == "" {
					return nil
				}
			case strings.HasPrefix(rule, "same:"):
				otherkey := strings.TrimPrefix(rule, "same:")
				if value != hold[otherkey] {
					return errors.New(key + " must match " + otherkey)
				}
			case strings.HasPrefix(rule, "exists:"):
				// Parse the table and column from the rule
				parts := strings.Split(strings.TrimPrefix(rule, "exists:"), ",")
				if len(parts) != 2 {
					return errors.New("invalid exists rule format, expected 'exists:table,column'")
				}
				table, column := parts[0], parts[1]
				if err := v.Exists(value, table, column); err != nil {
					return err
				}
			default:
				return errors.New("unknown validation rule: " + rule)
			}
		case func(interface{}) error: // Custom validation function
			if err := rule(value); err != nil {
				return err
			}
		default:
			return errors.New("invalid validation rule type")
		}
	}
	return nil
}

// Exists checks if a value exists in the specified table and column
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

func ValidateRequest(r *http.Request, inputs map[string][]interface{}, app *app.Application) (bool, map[string]string) {
	r.ParseForm()

	v := NewValidator(app)
	errors := make(map[string]string)

	hold := make(map[string]interface{})

	for key, _ := range inputs {
		value := r.FormValue(key)
		hold[key] = value
	}

	for key, rules := range inputs {
		value := r.FormValue(key)
		if err := v.ValidateInput(value, rules, key, hold); err != nil {
			errors[key] = err.Error()
		}
	}

	for index, errorMessage := range errors {
		errors[index] = helpers.BeautifyMessage(errorMessage)
	}

	return len(errors) == 0, errors
}
