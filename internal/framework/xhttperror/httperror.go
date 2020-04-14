package xhttperror

import (
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error with HTTP status code and error message
type XHTTPError interface {
	error
	// StatusCode returns the HTTP status code of the error
	StatusCode() int
}

type HttpCodeError struct {
	statusCode int
}

// NewHTTPStatusCodeError creates a new HttpError instance.
// to generate the message based on the status code.
func NewHTTPStatusCodeError(status int) XHTTPError {
	return &HttpCodeError{status}
}

func (e *HttpCodeError) Error() string {
	return http.StatusText(e.statusCode)
}

// StatusCode returns the HTTP status code.
func (e *HttpCodeError) StatusCode() int {
	return e.statusCode
}

type HttpJsonError struct {
	statusCode int
	Message    interface{}
}

// NewJsonHTTPError creates a new HttpError instance.
// to generate the message based on the status code.
func NewJsonHTTPError(status int, message interface{}) XHTTPError {
	return &HttpJsonError{statusCode: status, Message: message}
}

// StatusCode returns the HTTP status code.
func (e *HttpJsonError) StatusCode() int {
	return e.statusCode
}

func (e *HttpJsonError) Error() string {
	return fmt.Sprintf("%s: %v", http.StatusText(e.statusCode), e.Message)
}
