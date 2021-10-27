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
			name:  "uuid V1",
			regex: xregex.UUID,
			value: "72b08610-3735-11ec-8d3d-0242ac130003",
		},
		// Version 2 UUIDs are generated in the same way as version 1 UUIDs
		// but the low part of the timestamp (the time_low field) is replaced by a 32-bit integer
		// for brevity, reusing the V1 value with only the version bits adjusted
		{
			name:  "uuid V2",
			regex: xregex.UUID,
			value: "72b08610-3735-21ec-8d3d-0242ac130003",
		},
		{
			name:  "uuid V3",
			regex: xregex.UUID,
			value: "cf16fe52-3365-3a1f-8572-288d8d2aaa46",
		},
		{
			name:  "uuid V4",
			regex: xregex.UUID,
			value: "1314076b-ad82-48b1-8a3a-b74ab3b4af77",
		},
		{
			name:  "uuid V5",
			regex: xregex.UUID,
			value: "41072533-0430-54ea-9e2a-d21a311fa0b4",
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
