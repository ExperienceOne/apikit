package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"time"

	"io"

	"github.com/go-ozzo/ozzo-routing"
	log "github.com/sirupsen/logrus"
	"github.com/ExperienceOne/apikit/middleware"
)

const (
	DateTimeFormat = "2006-01-02T15:04:05"
)

func TestLogger(t *testing.T) {

	tests := []struct {
		name            string
		setup           func() (string, func())
		httpMethod      string
		requestEndpoint string
		requestBody     io.Reader
		contentType     string
		wantRequest     string
		wantResponse    string
	}{
		{
			name: "GET",
			setup: func() (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewLogger(DateTimeFormat, []string{}, func(entity middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", entity).Info("request handled")
				})).Handler)
				router.Get("/endpoint", func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"]}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodGet,
			requestEndpoint: "/endpoint",
			requestBody:     nil,
			contentType:     "application/json",
			wantRequest:     `"Request":{"Params":{},"Endpoint":"/endpoint","Method":"GET"}`,
			wantResponse:    `"Response":{"Body":{"number":6.13,"strings":["a","b"]},"responseCode":200}`,
		},
		{
			name: "GET_WithParameter",
			setup: func() (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewLogger(DateTimeFormat, []string{}, func(entity middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", entity).Info("request handled")
				})).Handler)
				router.Get(`/endpoint/something/<id:\d+><test>`, func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"]}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodGet,
			requestEndpoint: "/endpoint/something/1337?test=123456",
			requestBody:     nil,
			contentType:     "application/json",
			wantRequest:     `"Request":{"Params":{"test":["123456"]},"Endpoint":"/endpoint/something/1337","Method":"GET"}`,
			wantResponse:    `"Response":{"Body":{"number":6.13,"strings":["a","b"]},"responseCode":200}`,
		},
		{
			name: "POST",
			setup: func() (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewLogger(DateTimeFormat, []string{}, func(entity middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", entity).Info("request handled")
				})).Handler)
				router.Post("/endpoint", func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"]}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodPost,
			requestEndpoint: "/endpoint",
			requestBody:     bytes.NewBuffer([]byte(`{"float":6.13,"strings":["a","b"]}`)),
			contentType:     "application/json",
			wantRequest:     `"Request":{"Body":{"float":6.13,"strings":["a","b"]},"Params":{},"Endpoint":"/endpoint","Method":"POST"}`,
			wantResponse:    `"Response":{"Body":{"number":6.13,"strings":["a","b"]},"responseCode":200}`,
		},
		{
			name: "POST_FormUrlEncoded",
			setup: func() (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewLogger(DateTimeFormat, []string{}, func(entity middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", entity).Info("request handled")
				})).Handler)
				router.Post("/endpoint", func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"]}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodPost,
			requestEndpoint: "/endpoint",
			requestBody:     bytes.NewBufferString(url.Values{"test": {"test-string"}}.Encode()),
			contentType:     "application/x-www-form-urlencoded",
			wantRequest:     `"Request":{"Body":"test=test-string","Params":{},"Endpoint":"/endpoint","Method":"POST"}`,
			wantResponse:    `"Response":{"Body":{"number":6.13,"strings":["a","b"]},"responseCode":200}`,
		},
		{
			name: "POST_IgnoredPath",
			setup: func() (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewLogger(DateTimeFormat, []string{"/documents"}, func(entity middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", entity).Info("request handled")
				})).Handler)
				router.Post("/documents", func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"]}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodPost,
			requestEndpoint: "/documents",
			requestBody:     bytes.NewBuffer([]byte(`{"float":6.13,"strings":["a","b"]}`)),
			contentType:     "application/json",
			wantRequest:     `"Request":{"Params":{},"Endpoint":"/documents","Method":"POST"}`,
			wantResponse:    `"Response":{"responseCode":200}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buffer bytes.Buffer
			log.SetOutput(&buffer)
			log.SetFormatter(&log.JSONFormatter{})
			log.SetLevel(log.InfoLevel)

			serverUrl, closeFunc := tt.setup()
			defer closeFunc()

			req, err := http.NewRequest(tt.httpMethod, serverUrl+tt.requestEndpoint, tt.requestBody)
			if err != nil {
				t.Fatalf("http.NewRequest() unexpected error = %v", err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			client := &http.Client{Timeout: time.Second * 5}
			_, err = client.Do(req)
			if err != nil {
				t.Fatalf("http.Client.Do() unexpected error = %v", err)
			}

			if !strings.Contains(buffer.String(), tt.wantRequest) {
				t.Fatalf("expected buffer to contain: %s; got: %s", tt.wantRequest, buffer.String())
			}

			if !strings.Contains(buffer.String(), tt.wantResponse) {
				t.Fatalf("expected buffer to contain: %s; got: %s", tt.wantResponse, buffer.String())
			}
		})
	}
}

func TestMaskingLogger(t *testing.T) {

	tests := []struct {
		name            string
		setup           func(fieldsToMask []string) (string, func())
		httpMethod      string
		requestEndpoint string
		requestBody     io.Reader
		contentType     string
		fieldsToMask    []string
		wantRequest     string
		wantResponse    string
	}{
		{
			name: "POST",
			setup: func(fieldsToMask []string) (string, func()) {

				router := routing.New()
				router.Use(middleware.Log(middleware.NewMaskingLogger(DateTimeFormat, fieldsToMask, nil, []string{}, func(logEntry middleware.LogEntry, values ...interface{}) {
					log.WithField("entry", logEntry).Info("request handled")
				})).Handler)
				router.Post("/endpoint", func(c *routing.Context) error {
					c.Response.WriteHeader(http.StatusOK)
					return c.Write([]byte(`{"number":6.13,"strings":["a","b"],"name":"tester","location":{"name":"street"}}`))
				})

				ts := httptest.NewServer(router)
				return ts.URL, ts.Close
			},
			httpMethod:      http.MethodPost,
			requestEndpoint: "/endpoint",
			requestBody:     bytes.NewBuffer([]byte(`{"float":6.13,"strings":["a","b"],"username":"tester","location":{"name":"street"}}`)),
			contentType:     "application/json",
			fieldsToMask:    []string{"username", "name"},
			wantRequest:     `"Request":{"Body":{"float":6.13,"location":{"name":""},"strings":["a","b"],"username":""},"Params":{},"Endpoint":"/endpoint","Method":"POST"}`,
			wantResponse:    `"Response":{"Body":{"location":{"name":""},"name":"","number":6.13,"strings":["a","b"]},"responseCode":200}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buffer bytes.Buffer
			log.SetOutput(&buffer)
			log.SetFormatter(&log.JSONFormatter{})
			log.SetLevel(log.InfoLevel)

			serverUrl, closeFunc := tt.setup(tt.fieldsToMask)
			defer closeFunc()

			req, err := http.NewRequest(tt.httpMethod, serverUrl+tt.requestEndpoint, tt.requestBody)
			if err != nil {
				t.Fatalf("http.NewRequest() unexpected error = %v", err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			client := &http.Client{Timeout: time.Second * 5}
			_, err = client.Do(req)
			if err != nil {
				t.Fatalf("http.Client.Do() unexpected error = %v", err)
			}

			if !strings.Contains(buffer.String(), tt.wantRequest) {
				t.Fatalf("expected buffer to contain: %s; got: %s", tt.wantRequest, buffer.String())
			}

			if !strings.Contains(buffer.String(), tt.wantResponse) {
				t.Fatalf("expected buffer to contain: %s; got: %s", tt.wantResponse, buffer.String())
			}
		})
	}
}
