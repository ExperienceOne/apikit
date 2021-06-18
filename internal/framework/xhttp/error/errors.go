package error

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewNotSupportedContentType(statusCode int, message string) error {
	return &NotSupportedContentType{
		message:    message,
		statusCode: statusCode,
	}
}

type NotSupportedContentType struct {
	message    string
	statusCode int
}

func (e *NotSupportedContentType) Error() string {
	return fmt.Sprintf("error unsupported media type (%s)", e.message)
}

func (e *NotSupportedContentType) StatusCode() int {
	return e.statusCode
}

var NewRequestObjectIsNilError = errors.New("request object is nil")

func NewErrUnknownResponse(code int) *ErrUnknownResponse {
	return &ErrUnknownResponse{
		code: code,
	}
}

type ErrUnknownResponse struct {
	code int
}

func (err *ErrUnknownResponse) Error() string {
	return fmt.Sprintf("unknown response status code '%d'", err.code)
}

func NewErrOnUnknownResponseCode(message string) *ErrOnUnknownResponseCode {
	return &ErrOnUnknownResponseCode{
		Message: message,
	}
}

type ErrOnUnknownResponseCode struct {
	Message string
}

func (err *ErrOnUnknownResponseCode) Error() string {
	return fmt.Sprintf(err.Message)
}
