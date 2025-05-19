package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ErrorMessage(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))

		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s characters", err.Field(), err.Param()))

		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at most %s characters", err.Field(), err.Param()))

		case "gte":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be greater than or equal to %s", err.Field(), err.Param()))

		case "lte":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be less than or equal to %s", err.Field(), err.Param()))

		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be one of [%s]", err.Field(), err.Param()))

		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}
