package openapi_test

import (
	"testing"

	"github.com/ExperienceOne/apikit/generator/openapi"
)

func TestLoadDoc(t *testing.T) {

	tests := []struct {
		name     string
		filePath string
		want     error
	}{
		{
			name:     "YAML",
			filePath: "./spec.yaml",
		},
		{
			name:     "JSON",
			filePath: "./spec.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := openapi.NewOpenApiSpecFromFile(test.filePath)
			if err != nil {
				t.Errorf("could create document object for file (%v)", err)
			}

			b, err := doc.Spec.MarshalJSON()
			if err != nil {
				t.Error(err)
			}

			t.Log(string(b))
		})
	}
}
