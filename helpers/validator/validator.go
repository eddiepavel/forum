package validator

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)


type Validator struct{}


func NewValidator() *Validator {
	return &Validator{}
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
			case strings.HasPrefix(rule, "same:"):
				otherkey := strings.TrimPrefix(rule, "same:")
				if value != hold[otherkey] {
					return errors.New(key + " must match " + otherkey)
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


func ValidateRequest(r *http.Request, inputs map[string][]interface{}) (bool, map[string]string) {
	r.ParseForm()

	v := NewValidator()
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

	return len(errors) == 0, errors
}
