package error

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewNotSupportedContentType(statusCode int, message string) error {
	return &notSupportedContentType{
		message:    message,
		statusCode: statusCode,
	}
}

type notSupportedContentType struct {
	message    string
	statusCode int
}

func (e *notSupportedContentType) Error() string {
	return fmt.Sprintf("error unsupported media type (%s)", e.message)
}

func (e *notSupportedContentType) StatusCode() int {
	return e.statusCode
}

var NewRequestObjectIsNilError = errors.New("request object is nil")

func NewUnknownResponseError(code int) error {
	return fmt.Errorf("unknown response status code '%d'", code)
}
