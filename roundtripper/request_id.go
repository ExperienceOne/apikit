package roundtripper

import (
	"github.com/ExperienceOne/apikit/pkg/requestid"
	"net/http"
)

// RequestID generates a new request id foreach outgoing request
// or extracts a request ID from the context and then sets a request id header
func RequestID() RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return Func(func(req *http.Request) (*http.Response, error) {

			reqID := requestid.Get(req.Context())
			if reqID == "" {
				reqID = requestid.Generate()
			}

			req.Header.Set(requestid.RequestIdHttpHeader, reqID)

			return next.RoundTrip(req)
		})
	}
}
