package middleware

import (
	"context"
	"github.com/ExperienceOne/apikit/internal/framework/xserver"
	"github.com/ExperienceOne/apikit/pkg/requestid"
	"github.com/go-ozzo/ozzo-routing"
	"net/http"
)

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

	reqID := req.Header.Get(requestid.RequestIdHttpHeader)
	if reqID == "" {
		reqID = requestid.Generate()
	}

	return context.WithValue(ctx, requestid.RequestIdKey, reqID)
}
