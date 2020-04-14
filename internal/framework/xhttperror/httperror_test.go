package xhttperror_test

import (
	"net/http"
	"testing"

	"github.com/ExperienceOne/apikit/internal/framework/xhttperror"
)

func Test(t *testing.T) {
	tests := []struct {
		err  xhttperror.XHTTPError
		name string
	}{
		{
			name: "status code error",
			err:  xhttperror.NewHTTPStatusCodeError(http.StatusInternalServerError),
		},
		{
			name: "json error",
			err:  xhttperror.NewJsonHTTPError(http.StatusInternalServerError, "hell, yeah"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if message := test.err.Error(); message == "" {
				t.Error("message of error is empty")
			}
		})
	}
}
