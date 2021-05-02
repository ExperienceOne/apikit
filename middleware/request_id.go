package middleware

import (
	"context"
	"fmt"
	"github.com/ExperienceOne/apikit/internal/framework/xserver"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gofrs/uuid"
	"net/http"
	"sync/atomic"
)

// requestIdCxtKey key type used to store a request ID.
type requestIdCxtKey int

// RequestIdKey is the key to store an unique request ID in a request context.
const RequestIdKey requestIdCxtKey = 0

const RequestIdHttpHeader = "X-Request-ID"

// incrementable fallback ID in case of UUID generation failure
var requestIdFallback uint64

// RequestID creates a middleware which creates a request ID and store the ID in the request context.
func RequestID() xserver.Middleware {

	return xserver.Middleware{
		Handler: func(c *routing.Context) error {

			ctx := contextWithRequestID(c.Request.Context(), c.Request)
			c.Request = c.Request.WithContext(ctx)

			return c.Next()
		},
	}
}

func contextWithRequestID(ctx context.Context, req *http.Request) context.Context {

	reqID := req.Header.Get(RequestIdHttpHeader)
	if reqID == "" {
		u2, err := uuid.NewV4()
		if err != nil {
			newFallbackId := atomic.AddUint64(&requestIdFallback, 1)
			reqID = fmt.Sprint(newFallbackId)
		} else {
			reqID = u2.String()
		}
	}

	return context.WithValue(ctx, RequestIdKey, reqID)
}

// GetRequestID returns a request ID for the given context; otherwise an empty string
// will be returned if the context is nil or the request ID can not be found.
func GetRequestID(ctx context.Context) string {

	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(RequestIdKey).(string); ok {
		return reqID
	}

	return ""
}
