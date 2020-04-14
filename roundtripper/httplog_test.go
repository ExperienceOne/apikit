package roundtripper_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"reflect"

	"github.com/ExperienceOne/apikit/roundtripper"
)

func TestLog(t *testing.T) {

	want := "request-method = GET;response-status = 418;"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer ts.Close()

	var buffer bytes.Buffer

	logFunc := func(req *http.Request, resp *http.Response) {
		fmt.Fprintf(&buffer, "request-method = %s;", req.Method)
		fmt.Fprintf(&buffer, "response-status = %d;", resp.StatusCode)
	}

	client := roundtripper.Use(&http.Client{}, roundtripper.Logger(logFunc))
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	got := buffer.String()
	if got == "" {
		t.Fatal("expected buffer not to be empty")
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
