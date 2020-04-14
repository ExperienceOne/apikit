package xclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ExperienceOne/apikit/internal/framework/hooks"
)

// Opts contains hooks and an optional context object
type Opts struct {
	Hooks hooks.HooksClient
	Ctx   context.Context
}

// HttpClientWrapper imposes common Client API conventions on a set of resource paths.
// The baseURL is expected to point to an HTTP path that is the parent
// of one or more resources.  The server should return a decodable API resource
// object, or an Status object which contains information about the reason for
// any failure.
type HttpClientWrapper struct {
	// base is the root URL for all invocations of the client
	BaseURL string

	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	*http.Client
}

// NewHttpClientWrapper wraps http.Client to extend this client additional features
func NewHttpClientWrapper(client *http.Client, baseUrl string) *HttpClientWrapper {
	return &HttpClientWrapper{
		Client:  client,
		BaseURL: baseUrl,
	}
}

type NewRequest func(string, io.Reader) (*http.Request, error)

// Verb begins a request with a http method verb.
func (c *HttpClientWrapper) Verb(verb string) NewRequest {
	baseURL := c.BaseURL
	return func(endpoint string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequest(verb, baseURL+endpoint, body)
		if err != nil {
			return nil, err
		}
		return req, err
	}
}

// Get begins a GET request. Short for c.Verb("GET").
func (c *HttpClientWrapper) Get() NewRequest {
	return c.Verb(http.MethodGet)
}

func (c *HttpClientWrapper) Into(body io.ReadCloser, r interface{}) error {
	err := json.NewDecoder(body).Decode(&r)
	if err != nil {
		return err
	}
	defer body.Close()
	return nil
}
