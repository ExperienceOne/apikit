package generator_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ExperienceOne/apikit/generator"
	"github.com/ExperienceOne/apikit/generator/openapi"
)

func Test(t *testing.T) {

	tests := []struct {
		name     string
		files    []string
		generate func(path, pkg string) error
		count    int
	}{
		{
			name: "client generator",
			generate: func(path, pkg string) error {

				spec, err := openapi.NewOpenApiSpecFromFile("../tests/data/swagger.yaml")
				if err != nil {
					return err
				}

				return generator.NewGoClientAPIGenerator(spec).Generate(path, pkg, false)
			},
			files: []string{
				"types.go",
				"framework.go",
				"client.go",
			},
			count: 3,
		},
		{
			name: "server generator",
			generate: func(path, pkg string) error {

				spec, err := openapi.NewOpenApiSpecFromFile("../tests/data/swagger.yaml")
				if err != nil {
					return err
				}

				return generator.NewGoServerAPIGenerator(spec).Generate(path, pkg, false)
			},
			files: []string{
				"types.go",
				"framework.go",
				"server.go",
			},
			count: 3,
		},
		{
			name: "api generator",
			generate: func(path, pkg string) error {

				spec, err := openapi.NewOpenApiSpecFromFile("../tests/data/swagger.yaml")
				if err != nil {
					return err
				}

				return generator.NewGoAPIGenerator(spec).Generate(path, pkg, true)
			},
			files: []string{
				"types.go",
				"framework.go",
				"server.go",
				"client.go",
			},
			count: 4,
		},
	}

	tmpDir := os.TempDir()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testDir := filepath.Join(tmpDir, "test"+strconv.Itoa(rand.Int()))

			if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
				t.Errorf("could not create directory")
			}

			if err := test.generate(testDir, "generic"); err != nil {
				t.Errorf("could not generate files (%v)", err)
				return
			}

			var got int
			err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					t.Log(path)

					if containsFile(test.files, info.Name()) {
						got++
					}
				}
				return nil
			})

			if err != nil {
				t.Error(err)
				return
			}

			if got != test.count {
				t.Errorf("count is bad, got=%d", got)
			}

		})
	}
}

func containsFile(files []string, src string) bool {
	for _, files := range files {
		if files == src {
			return true
		}
	}
	return false
}
