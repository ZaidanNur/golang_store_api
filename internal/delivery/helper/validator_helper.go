package helper

import (
	"github.com/go-playground/validator/v10"
)

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "Must be greater than " + fe.Param()
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + fe.Param() + " characters"
	case "max":
		return "Must be at most " + fe.Param() + " characters"
	default:
		return "Invalid value"
	}
}
