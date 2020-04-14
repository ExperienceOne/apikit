package roundtripper

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Dump creates a round tripper to dump outgoing HTTP requests and incoming responses
// of a http.Client for analysing purpose.
func Dump(logFunc func(a ...interface{})) RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return Func(func(req *http.Request) (*http.Response, error) {

			resp, err := next.RoundTrip(req)
			if err != nil {
				return nil, err
			}

			dumpResponse, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return nil, err
			}

			dumpRequest, err := httputil.DumpRequest(req, true)
			if err != nil {
				return nil, err
			}

			logFunc(fmt.Sprintf("HTTP request:\n %v \n", string(dumpRequest)))
			logFunc(fmt.Sprintf("HTTP response:\n %v \n", string(dumpResponse)))

			return resp, err
		})
	}
}
