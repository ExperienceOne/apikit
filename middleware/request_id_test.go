package middleware_test

import (
	"bytes"
	"github.com/ExperienceOne/apikit/pkg/requestid"
	"testing"

	"net/http/httptest"

	"net/http"

	"github.com/ExperienceOne/apikit/middleware"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {

	var buf bytes.Buffer

	router := routing.New()
	router.Use(middleware.RequestID().Handler)
	router.Get("/", func(c *routing.Context) error {

		reqID := requestid.Get(c.Request.Context())
		buf.WriteString(reqID)

		return c.Write("done")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := http.Get(ts.URL + "/")
	require.Nil(t, err, "failed to send get request")

	t.Log(buf.String())

	require.NotEqual(t, buf.String(), "")
}
