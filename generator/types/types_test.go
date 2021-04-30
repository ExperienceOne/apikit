package types_test

import (
	"testing"

	"github.com/ExperienceOne/apikit/generator/types"
	"github.com/go-openapi/spec"
)

type ExpectedSchemaType struct {
	Composit types.Composit
	Name     string
	Type     string
	Required bool
	Elements []string
}

func TestGoSchemaType(t *testing.T) {

	testCases := []struct {
		name  string
		input struct {
			schema   *spec.Schema
			required bool
		}
		Output ExpectedSchemaType
	}{
		{
			name: "string (non pointer)",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"string"},
					},
				},
				required: true,
			},
			Output: ExpectedSchemaType{Type: "string", Required: true},
		},
		{
			name: "[]string (non pointer)",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"array"},
						Items: &spec.SchemaOrArray{
							Schema: &spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
						},
					},
				},
				required: true,
			},
			Output: ExpectedSchemaType{Composit: types.Array, Type: "string", Required: true},
		},
		{
			name: "boolean (non pointer)",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"boolean"},
					},
				},
				required: true,
			},
			Output: ExpectedSchemaType{Type: "bool", Required: true},
		},
		{
			name: "object (non pointer)",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
					},
				},
				required: true,
			},
			Output: ExpectedSchemaType{Composit: types.Object, Name: "Object1", Required: true},
		},
		{
			name: "pointer to a string",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"string"},
					},
				},
				required: false,
			},
			Output: ExpectedSchemaType{Type: "string", Required: false},
		},
		{
			name: "pointer to boolean",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"boolean"},
					},
				},
				required: false,
			},
			Output: ExpectedSchemaType{Type: "bool", Required: false},
		},
		{
			name: "pointer to object",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
					},
				},
				required: false,
			},
			Output: ExpectedSchemaType{Composit: types.Object, Name: "Object2", Required: false},
		},
		{
			name: "[]string",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"array"},
						Items: &spec.SchemaOrArray{
							Schema: &spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
						},
					},
				},
				required: false,
			},
			Output: ExpectedSchemaType{Composit: types.Array, Type: "string", Required: false},
		},
		{
			name: "object with required field",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:     []string{"object"},
						Required: []string{"Name"},
						Properties: map[string]spec.Schema{
							"Name": spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
						},
					},
				},
				required: true,
			},
			Output: ExpectedSchemaType{Composit: types.Object, Name: "Object3", Required: true, Elements: []string{"Name"}},
		},
		{
			name: "object with not required field",
			input: struct {
				schema   *spec.Schema
				required bool
			}{
				schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
						Properties: map[string]spec.Schema{
							"Name": spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
						},
					},
				},
				required: false,
			},
			Output: ExpectedSchemaType{Composit: types.Object, Name: "Object4", Required: false, Elements: []string{"Name"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			typ, err := types.FromSchema("", testCase.input.schema, testCase.input.required, nil)
			if err != nil {
				t.Fatal(err)
			}

			if typ.Composit != testCase.Output.Composit {
				t.Errorf(`unexpected composite type (actual: "%d", expected: "%d")`, typ.Composit, testCase.Output.Composit)
			}

			if typ.Name != testCase.Output.Name {
				t.Errorf(`unexpected type name (actual: "%s", expected: "%s")`, typ.Name, testCase.Output.Name)
			}

			if typ.Type != testCase.Output.Type {
				t.Errorf(`unexpected type (actual: "%s", expected: "%s")`, typ.Type, testCase.Output.Type)
			}

			if typ.Required != testCase.Output.Required {
				t.Errorf(`unexpected required value (actual: "%t", expected: "%t")`, typ.Required, testCase.Output.Required)
			}

			for _, element := range testCase.Output.Elements {
				found := false
				for _, property := range typ.Elements {
					if property.Name == element {
						found = true
						break
					}
				}
				if !found {
					t.Errorf(`missing property, expected: "%s"`, element)
				}
			}
		})
	}
}
