package validation

import (
	"fmt"

	"github.com/go-playground/validator"
)

type ValidationErrorsObject struct {
	Message string                  `json:"message"`
	Errors  []ValidationErrorObject `json:"errors"`
}

type ValidationErrorObject struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Code    string `json:"code"`
}

func NewValidation() *Validator {
	return &Validator{
		validator.New(),
	}
}

type Validator struct {
	*validator.Validate
}

// validates an incoming request and returns nil, nil if request is valid
// if an error occurred, it is an internal validation error and not a validation issue of the request
func (v *Validator) ValidateRequest(request interface{}) (*ValidationErrorsObject, error) {

	if err := v.Struct(request); err != nil {

		validationErrors := new(ValidationErrorsObject)
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, err
		}

		validationErrors.Message = "validation failed"
		for _, err := range errors {

			errorCode := fmt.Sprintf("invalid-%s", err.Tag())
			if err.Tag() == "required" {
				errorCode = err.Tag()
			}

			validationError := ValidationErrorObject{
				Message: fmt.Sprint(err),
				Field:   err.Field(),
				Code:    errorCode,
			}
			validationErrors.Errors = append(validationErrors.Errors, validationError)
		}

		return validationErrors, nil
	}

	return nil, nil
}
