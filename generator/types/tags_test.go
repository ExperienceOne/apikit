package types

import (
	"testing"

	"github.com/go-openapi/spec"
)

type RestrictionCase struct {
	Name       string
	makeSchema func() *spec.Schema
	want       string
}

func TestStringRestrictionTag(t *testing.T) {

	tests := []RestrictionCase{
		{
			Name: "number range",
			makeSchema: func() *spec.Schema {
				min := int64(1)
				max := int64(3)
				schema := &spec.Schema{}
				schema.MaxLength = &max
				schema.MinLength = &min
				return schema
			},
			want: `min=1,max=3`,
		},
		{
			Name: "number max",
			makeSchema: func() *spec.Schema {
				max := int64(3)
				schema := &spec.Schema{}
				schema.MaxLength = &max
				return schema
			},
			want: `max=3`,
		},
		{
			Name: "number min",
			makeSchema: func() *spec.Schema {
				min := int64(1)
				schema := &spec.Schema{}
				schema.MinLength = &min
				return schema
			},
			want: `min=1`,
		},
		{
			Name: "email",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				schema.Format = "email"
				return schema
			},
			want: `email`,
		},
		{
			Name: "pattern regex and min and max",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := int64(1)
				max := int64(3)
				schema.MaxLength = &max
				schema.MinLength = &min
				return schema
			},
			want: `min=1,max=3`,
		},
		{
			Name: "pattern regex",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := int64(3)
				schema.MinLength = &min
				return schema
			},
			want: `min=3`,
		},
		{
			Name: "pattern regex and max",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				max := int64(3)
				schema.MaxLength = &max
				return schema
			},
			want: `max=3`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			schema := test.makeSchema()
			got := generateStringRestriction(schema.MaxLength, schema.MinLength, schema.Format)
			tag := generateValidationTag(true, got)
			if tag != test.want {
				t.Errorf("restriction is bad (want: %s got: %s)", test.want, tag)
			}
			t.Logf("tag: %s", got)
		})
	}
}

func TestIntegerRestrictionTag(t *testing.T) {

	tests := []RestrictionCase{
		{
			Name: "Maximum",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				max := float64(3)
				schema.Maximum = &max
				return schema
			},
			want: `max=3`,
		},
		{
			Name: "Minimum",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := float64(3)
				schema.Minimum = &min
				return schema
			},
			want: `min=3`,
		},
		{
			Name: "Minimum and maximum",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := float64(1)
				schema.Minimum = &min
				max := float64(3)
				schema.Maximum = &max
				return schema
			},
			want: `min=1,max=3`,
		},
		{
			Name: "Minimum exclusive",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := float64(1)
				schema.Minimum = &min
				schema.ExclusiveMinimum = true
				return schema
			},
			want: `min=2`,
		},
		{
			Name: "Maximum exclusive",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				max := float64(1)
				schema.Maximum = &max
				schema.ExclusiveMaximum = true
				return schema
			},
			want: `max=0`,
		},
		{
			Name: "Minimum and maximum (exclusive)",
			makeSchema: func() *spec.Schema {
				schema := &spec.Schema{}
				min := float64(1)
				schema.Minimum = &min
				schema.ExclusiveMinimum = true
				max := float64(3)
				schema.Maximum = &max
				schema.ExclusiveMaximum = true
				return schema
			},
			want: `min=2,max=2`,
		},
	}

	for _, test := range tests {
		schema := test.makeSchema()
		got := generateIntegerRestriction(schema.Minimum, schema.Maximum, schema.ExclusiveMinimum, schema.ExclusiveMaximum)
		tag := generateValidationTag(true, got)
		if tag != test.want {
			t.Errorf("restriction is bad (want: %s got: %s)", test.want, tag)
		}
		t.Logf("tag: %s", got)
	}
}
