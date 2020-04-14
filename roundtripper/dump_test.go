package roundtripper_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ExperienceOne/apikit/roundtripper"
)

func TestDump(t *testing.T) {

	var buff bytes.Buffer
	logger := func(a ...interface{}) {
		buff.WriteString(a[0].(string))
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		if _, err := w.Write([]byte("1")); err != nil {
			t.Fatalf("failed to write test byte on response writer: %s", err)
		}
	}))
	defer ts.Close()

	client := roundtripper.Use(&http.Client{}, roundtripper.Dump(logger))
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("response code is bad, got=%d", http.StatusTeapot)
	}

	if !strings.Contains(buff.String(), "418") {
		t.Error("status code is missing")
	}

	if !strings.Contains(buff.String(), "HTTP response:") {
		t.Error("HTTP response is missing")
	}

	if !strings.Contains(buff.String(), "HTTP request:") {
		t.Error("HTTP request is missing")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "1" {
		t.Errorf("body is bad, got=%v", 1)
	}
}
