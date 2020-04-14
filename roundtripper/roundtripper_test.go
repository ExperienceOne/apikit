package roundtripper_test

import (
	"bytes"
	"testing"

	"net/http"
	"net/http/httptest"

	"fmt"

	"reflect"

	"github.com/ExperienceOne/apikit/roundtripper"
)

func testRoundTripper(count int, buffer *bytes.Buffer) roundtripper.RoundTripper {

	return func(next http.RoundTripper) http.RoundTripper {
		return roundtripper.Func(func(r *http.Request) (*http.Response, error) {

			fmt.Fprint(buffer, count)
			return next.RoundTrip(r)
		})
	}
}

func TestUse(t *testing.T) {

	want := "210"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	var buffer bytes.Buffer

	numTripper := 3
	tripper := make([]roundtripper.RoundTripper, numTripper)
	for i := 0; i < numTripper; i++ {
		tripper[i] = testRoundTripper(i, &buffer)
	}

	client := roundtripper.Use(&http.Client{}, tripper...)
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	got := buffer.String()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
