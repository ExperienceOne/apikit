package xcontext

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type contextKey int

const (
	RequestHeaderKey contextKey = 1 + iota
)

// HttpContext is a context bucket for http related objects like headers
type HttpContext interface {
	GetHTTPRequestHeaders() (http.Header, bool)
}

func CreateHttpContext(header http.Header) context.Context {
	ctx := context.Background()
	ctx = hTTPRequestHeaders(ctx, header)
	return ctx
}

// GetHTTPRequestHeaders injects header map into context
func hTTPRequestHeaders(ctx context.Context, header http.Header) context.Context {
	return context.WithValue(ctx, RequestHeaderKey, header)
}

// NewHttpContextWrapper creates a wrapper for a context object
func NewHttpContextWrapper(ctx context.Context) HttpContext {
	if ctx == nil {
		ctx = context.Background()
	}
	return &httpContext{
		ctx,
	}
}

type httpContext struct {
	context.Context
}

// GetHTTPRequestHeaders retrieves header map form context
func (c *httpContext) GetHTTPRequestHeaders() (http.Header, bool) {
	header, ok := c.Value(RequestHeaderKey).(http.Header)
	return header, ok
}

// SetRequestHeadersFromContext iterates over headers in context and add adds all values to the request headers map
func SetRequestHeadersFromContext(httpContext HttpContext, header http.Header) error {

	if httpContext == nil {
		return nil
	}

	headersFromContext, ok := httpContext.GetHTTPRequestHeaders()
	if !ok {
		return nil
	}

	// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
	for key, values := range headersFromContext {
		if _, exists := header[key]; exists {
			return errors.New("header from context overwrites header in request object")
		}

		header[key] = values
	}
	return nil
}
