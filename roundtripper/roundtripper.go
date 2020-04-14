package roundtripper

import "net/http"

// RoundTripper represents an equivalent to a server middleware for a http client.
type RoundTripper func(next http.RoundTripper) http.RoundTripper

// Func transforms a function into a http.RoundTripper.
type Func func(req *http.Request) (*http.Response, error)

func (f Func) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Use appends RoundTripper to the transport layer of a http.Client.
func Use(client *http.Client, tripper ...RoundTripper) *http.Client {

	if len(tripper) == 0 {
		return client
	}

	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	for _, t := range tripper {
		transport = t(transport)
	}

	updatedClient := *client
	updatedClient.Transport = transport

	return &updatedClient
}
