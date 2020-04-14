package xcontext_test

import (
	"github.com/ExperienceOne/apikit/internal/framework/xhttp/xcontext"
	"net/http"
	"testing"
)

func TestSetRequestHeadersFromContext(t *testing.T) {
	contextApiKey := "x-context-api-key"
	contextApiKeyValue := "test"
	headerInContext := make(http.Header)
	headerInContext.Set(contextApiKey, contextApiKeyValue)
	headerInRequest := make(http.Header)
	ctx := xcontext.CreateHttpContext(headerInContext)
	object := xcontext.NewHttpContextWrapper(ctx)
	if err := xcontext.SetRequestHeadersFromContext(object, headerInRequest); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	value := headerInRequest.Get(contextApiKey)
	if value != contextApiKeyValue {
		t.Errorf("unexpected error (actual: '%s', expected: '%s')", value, contextApiKeyValue)
	}
}
