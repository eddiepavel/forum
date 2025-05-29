package validator

import (
	"errors"
	"forum-app/app"
	"forum-app/helpers"
	"net/http"
	"strings"
)

type Validator struct {
	app *app.Application
}

func NewValidator(app *app.Application) *Validator {
	return &Validator{app: app}
}

// ValidateInput validates a value against a set of rules and updates the hold map.
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
			case rule == "login_attempt":
				email, emailExists := hold["email"].(string)
				password, passwordExists := hold["password"].(string)
				if !emailExists || !passwordExists {
					return errors.New("email and password are required for login attempt validation")
				}
				if err := v.ValidateLoginAttempt(email, password); err != nil {
					return err
				}
			case rule == "password":
				if password, ok := value.(string); ok {
					if err := ValidatePassword(password); err != nil {
						return err
					}
				} else {
					return errors.New(key + " must be a valid string")
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

// ValidateRequest validates HTTP request inputs based on provided rules and returns errors if any.
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

	// Remove duplicate error messages
	uniqueErrors := make(map[string]bool)
	for key, errorMessage := range errors {
		if uniqueErrors[errorMessage] {
			delete(errors, key) // Remove duplicate error
		} else {
			uniqueErrors[errorMessage] = true
		}
	}

	return len(errors) == 0, errors
}
