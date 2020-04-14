package roundtripper

import "net/http"

// Log creates a round tripper to log outgoing HTTP requests and incoming responses
// of a http.Client.
func Logger(logFunc func(*http.Request, *http.Response)) RoundTripper {

	return func(next http.RoundTripper) http.RoundTripper {
		return Func(func(req *http.Request) (*http.Response, error) {

			resp, err := next.RoundTrip(req)
			logFunc(req, resp)

			return resp, err
		})
	}
}
