package xregex_test

import (
	"regexp"
	"testing"

	"github.com/ExperienceOne/apikit/generator/xregex"
)

func Test(t *testing.T) {
	tests := []struct {
		name  string
		regex string
		value string
	}{
		{
			name:  "URL http",
			regex: xregex.URL,
			value: "http://example.de",
		},
		{
			name:  "URL https",
			regex: xregex.URL,
			value: "https://example.de",
		},
		{
			name:  "URL https and sub path",
			regex: xregex.URL,
			value: "https://example.de/test/test",
		},
		{
			name:  "URL https and sub path and query",
			regex: xregex.URL,
			value: "https://example.de/test/test?query=9",
		},
		{
			name:  "uuid",
			regex: xregex.UUID,
			value: "1314076b-ad82-48b1-8a3a-b74ab3b4af77",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			regex, err := regexp.Compile(test.regex)
			if err != nil {
				t.Fatalf("couldn't compile regex (%v)", err)
			}
			if !regex.MatchString(test.value) {
				t.Error("value didn't machted regex")
			}
		})
	}
}
