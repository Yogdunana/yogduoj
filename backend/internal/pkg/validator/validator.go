package validator

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// ValidatePassword checks password strength requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 128 {
		return ErrPasswordTooLong
	}

	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return ErrPasswordWeak
	}

	return nil
}

// ValidateUsername checks username format.
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	if len(username) > 32 {
		return ErrUsernameTooLong
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return ErrUsernameInvalid
	}

	return nil
}

// ValidateEmail checks email format.
func ValidateEmail(email string) error {
	if email == "" {
		return nil
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return ErrEmailInvalid
	}
	return nil
}

// GetStructTags returns the json tag names for a struct's fields.
func GetStructTags(v interface{}) []string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	var tags []string
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			parts := strings.Split(tag, ",")
			tags = append(tags, parts[0])
		}
	}
	return tags
}

// BindJSON binds JSON request body and returns validation errors.
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

// Various validation errors
var (
	ErrPasswordTooShort = &ValidationError{Field: "password", Message: "password must be at least 8 characters"}
	ErrPasswordTooLong  = &ValidationError{Field: "password", Message: "password must be at most 128 characters"}
	ErrPasswordWeak     = &ValidationError{Field: "password", Message: "password must contain uppercase, lowercase, and digit"}
	ErrUsernameTooShort = &ValidationError{Field: "username", Message: "username must be at least 3 characters"}
	ErrUsernameTooLong  = &ValidationError{Field: "username", Message: "username must be at most 32 characters"}
	ErrUsernameInvalid  = &ValidationError{Field: "username", Message: "username can only contain letters, numbers, underscores, and hyphens"}
	ErrEmailInvalid     = &ValidationError{Field: "email", Message: "invalid email format"}
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
