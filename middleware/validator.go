package middleware

import (
	"github.com/go-playground/validator"
)

var Validator = &CustomValidator{validator: validator.New()}

// CustomValidator struktur
type CustomValidator struct {
	validator *validator.Validate
}

// Custom validation function for Echo
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
