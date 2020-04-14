package xhttp_test

import (
	"testing"

	"github.com/ExperienceOne/apikit/internal/framework/xhttp"
)

func TestExtractContentType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		ouput string
	}{
		{
			name:  "empty header",
			input: "",
			ouput: "",
		},
		{
			name:  "header with ;",
			input: "text/plain; charset=utf-8",
			ouput: xhttp.ContentTypeTextPlain,
		},
		{
			name:  "header without ;",
			input: xhttp.ContentTypeApplicationJson,
			ouput: xhttp.ContentTypeApplicationJson,
		},
		{
			name:  "header with hal",
			input: xhttp.ContentTypeApplicationHalJson,
			ouput: xhttp.ContentTypeApplicationHalJson,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := xhttp.ExtractContentType(test.input)
			if test.ouput != output {
				t.Errorf(`unexpected header value (actual:"%s", expected: "%s")`, output, test.ouput)
			}
		})
	}
}

func TestContentTypeInList(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			list   []string
			search string
		}
		ouput bool
	}{
		{
			name: "in list",
			input: struct {
				list   []string
				search string
			}{
				list:   []string{"contentType1", "contentType2"},
				search: "contentType2",
			},
			ouput: true,
		},
		{
			name: "not in list",
			input: struct {
				list   []string
				search string
			}{
				list:   []string{"contentType1", "contentType2"},
				search: "contentType3",
			},
			ouput: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := xhttp.ContentTypeInList(test.input.list, test.input.search)
			if test.ouput != output {
				t.Errorf(`unexpected search result (actual:"%t", expected: "%t")`, output, test.ouput)
			}
		})
	}
}
