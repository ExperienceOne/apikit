package error_test

import (
	"github.com/stretchr/testify/assert"
	httperror "github.com/ExperienceOne/apikit/internal/framework/xhttp/error"
	"net/http"
	"testing"
)

func TestNotSupportedContentType(t *testing.T) {

	tests := []struct {
		name    string
		error   error
		message string
	}{
		{
			name:    "content type not supported",
			error:   httperror.NewNotSupportedContentType(http.StatusUnsupportedMediaType, "hello, world"),
			message: "error unsupported media type (hello, world)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.error == nil {
				t.Error("error is nil")
				return
			}
			assert.Equal(t, test.error.Error(), test.message)
		})
	}
}
