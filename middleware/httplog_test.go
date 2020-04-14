package middleware_test

import (
	"bytes"
	"net/http"
	"testing"

	"fmt"

	"io/ioutil"

	"net/http/httptest"

	"time"

	"net/url"

	"github.com/ExperienceOne/apikit/middleware"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogHandler(t *testing.T) {

	var buf bytes.Buffer
	logFunc := func(r *http.Request, w *middleware.LogResponseWriter, elapsed time.Duration) {
		req, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(&buf, "request: [%s %s %s] response: [%d %s] duration: [%f]", r.Method, r.URL.String(), string(req), w.Status, string(w.Body), float64(elapsed.Nanoseconds())/1e6)
	}

	router := routing.New()
	router.Use(middleware.LogHandler(logFunc))
	router.Get("/get", func(c *routing.Context) error {
		c.Response.WriteHeader(200)
		return c.Write("passing the test")
	})
	router.Post("/post", func(c *routing.Context) error {
		c.Response.WriteHeader(200)
		return c.Write("passing the test")
	})
	router.Post("/urlencoded", func(c *routing.Context) error {
		c.Response.WriteHeader(200)
		return c.Write("passing the test")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := http.Get(ts.URL + "/get")
	require.Nil(t, err, "failed to send get request %s")

	assert.Contains(t, buf.String(), "request: [GET /get ]")
	assert.Contains(t, buf.String(), "response: [200 passing the test]")

	buf.Reset()
	_, err = http.Post(ts.URL+"/post", "application/json", bytes.NewBuffer([]byte(`{"name": "test"}`)))
	require.Nil(t, err, "failed to send post request %s")

	assert.Contains(t, buf.String(), `request: [POST /post {"name": "test"}]`)
	assert.Contains(t, buf.String(), "response: [200 passing the test]")

	buf.Reset()
	form := url.Values{"test": {"test-string"}}
	_, err = http.Post(ts.URL+"/urlencoded", "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	require.Nil(t, err, "failed to send post request %s")

	assert.Contains(t, buf.String(), `request: [POST /urlencoded test=test-string]`)
	assert.Contains(t, buf.String(), "response: [200 passing the test]")
}
