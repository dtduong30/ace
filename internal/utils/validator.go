package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// ValidateStruct validates a struct and returns a list of user-friendly error messages
func ValidateStruct(data interface{}) []string {
	err := Validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		var msg string
		switch err.Tag() {
		case "required":
			msg = fmt.Sprintf("%s là bắt buộc", err.Field())
		case "min":
			msg = fmt.Sprintf("%s phải có giá trị tối thiểu là %s", err.Field(), err.Param())
		case "max":
			msg = fmt.Sprintf("%s không được vượt quá %s", err.Field(), err.Param())
		case "oneof":
			msg = fmt.Sprintf("%s phải là một trong: %s", err.Field(), err.Param())
		default:
			msg = fmt.Sprintf("%s không hợp lệ (%s)", err.Field(), err.Tag())
		}
		errors = append(errors, msg)
	}
	return errors
}
