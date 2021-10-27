package roundtripper

import (
	"bytes"
	"github.com/ExperienceOne/apikit/pkg/requestid"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {

	var buf bytes.Buffer

	router := routing.New()
	router.Get("/", func(c *routing.Context) error {

		reqID := c.Request.Header.Get(requestid.RequestIdHttpHeader)

		t.Log(reqID)

		buf.WriteString(reqID)

		return c.Write("done")
	})

	httpClient := new(http.Client)
	httpClient = Use(httpClient, RequestID())

	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := httpClient.Get(ts.URL + "/")
	require.Nil(t, err, "failed to send get request")

	require.NotEqual(t, buf.String(), "")
}
