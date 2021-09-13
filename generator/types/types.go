package types

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ExperienceOne/apikit/generator/file"
	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/ExperienceOne/apikit/generator/xregex"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	BooleanType = "bool"
	StringType  = "string"
	Float64Type = "float64"
	Float32Type = "float32"
	Int64Type   = "int64"
	IntType     = "int"
	Int32Type   = "int32"
	ByteType    = "byte"
)

func ConvertSimpleType(typ string, format string) string {

	var goType string

	if typ == "string" {
		goType = StringType
	} else if typ == "boolean" {
		goType = BooleanType
	} else if typ == "date" || typ == "date-time" || typ == "password" || typ == "byte" || typ == "binary" {
		goType = StringType
	} else if typ == "number" || typ == "float" || typ == "double" {
		goType = Float64Type
	} else if typ == "integer" && format == "int64" {
		goType = Int64Type
	} else if typ == "integer" && format == "int32" {
		goType = Int32Type
	} else if typ == "integer" {
		goType = Int64Type
	}

	return goType
}

// MatchTypes verifies types falls into the same primitive type category
func MatchTypes(param interface{}, typ, format string) bool {
	if typ == "number" || typ == "float" || typ == "double" {
		if _, ok := param.(float64); ok {
			return true
		}
	} else if typ == "integer" && format == "int64" {
		v, ok := param.(float64)
		if ok && float64(int64(v)) == v {
			return true
		}
	} else if typ == "integer" && format == "int32" {
		v, ok := param.(float64)
		if ok && float64(int32(v)) == v {
			return true
		}
	} else if typ == "integer" {
		v, ok := param.(float64)
		if ok && float64(int(v)) == v {
			return true
		}
	}
	return false
}

type Composit int

const (
	Simple Composit = iota
	Enum
	Object
	Map
	Array
	File
)

type Element struct {
	Name       string
	Type       *Type
	Serialized string
}

type Type struct {
	Composit   Composit
	Name       string
	Type       string
	Validation string
	Nestable   bool
	Required   bool
	Elements   []Element
	Validator  *RegexValidator
}

func New(composit Composit, name, typ string, required, nestable bool, tags ...string) *Type {

	return &Type{
		Composit:   composit,
		Name:       name,
		Type:       typ,
		Required:   required,
		Nestable:   nestable,
		Validation: generateValidationTag(required, tags),
	}
}

func NewObject(name string, required bool) *Type {

	return New(Object, name, "", required, true)
}

func NewArray(name string, required bool, elementType *Type) *Type {

	var tags []string
	if elementType.Nestable {
		tags = append(tags, "dive")
	}

	arrayType := elementType.Name
	if arrayType == "" {
		arrayType = elementType.Type
	}

	typ := New(Array, name, arrayType, required, false, tags...)
	typ.AddElement("", elementType, "")
	return typ
}

func NewMap(name string, required bool, elementType *Type) *Type {

	var tags []string
	if elementType.Nestable {
		tags = append(tags, "dive")
	}

	mapType := elementType.Name
	if mapType == "" {
		mapType = elementType.Type
	}

	typ := New(Map, name, mapType, required, true, tags...)
	typ.AddElement("", elementType, "")
	return typ
}

func NewEnum(name string, required bool, values []interface{}) (*Type, error) {

	if len(values) == 0 {
		return nil, errors.New("enum is empty")
	}

	var valueType string
	switch typeName := values[0].(type) {
	case bool:
		valueType = BooleanType
	case float32:
		valueType = Float32Type
	case float64:
		valueType = Float64Type
	case int:
		valueType = IntType
	case int32:
		valueType = Int32Type
	case int64:
		valueType = Int64Type
	case string:
		valueType = StringType
	case byte:
		valueType = ByteType
	default:
		return nil, errors.Errorf("enum type is not supported (%s)", typeName)
	}

	enumName := name
	typ := New(Enum, enumName, valueType, required, false)
	for _, value := range values {
		name := strings.Title(identifier.MakeIdentifier(fmt.Sprintf("%s_%#v", strings.Title(name), value)))
		typ.AddElement(name, &Type{Name: fmt.Sprintf("%#v", value), Type: enumName}, "")
	}

	return typ, nil
}

func NewFile(name string, required bool) *Type {

	return New(File, name, "io.ReadCloser", required, false)
}

func makeStringValidations(name, format string, validations *spec.CommonValidations) ([]string, *RegexValidator, error) {

	//note that an empty string is a valid string unless minLength or pattern is specified.

	var err error
	var tags []string
	var validator *RegexValidator

	if format == openapi.UUID {
		tags = generateStringRestriction(validations.MaxLength, validations.MinLength, format)
		validator, err = generateRegexRestriction(name, xregex.UUID)
	} else if format == openapi.URL {
		tags = generateStringRestriction(validations.MaxLength, validations.MinLength, format)
		validator, err = generateRegexRestriction(name, xregex.URL)
	} else {
		tags = generateStringRestriction(validations.MaxLength, validations.MinLength, format)
		validator, err = generateRegexRestriction(name, validations.Pattern)
	}

	if err != nil {
		return nil, nil, err
	}
	if validator != nil {
		tags = append(tags, validator.Tag)
	}

	return tags, validator, nil
}

func FromSimpleSchema(name string, schema *spec.SimpleSchema, required bool, validations *spec.CommonValidations) (*Type, error) {

	if schema == nil {
		return nil, errors.Errorf("'%s' has no schema", name)
	}

	var typ, format string
	if schema.Items == nil {
		typ = schema.Type
		format = schema.Format
	} else {
		typ = schema.Items.Type
		format = schema.Items.Format
	}

	goType := ConvertSimpleType(typ, format)
	if goType == "" {
		return nil, errors.Errorf("unknown simple type '%s' / '%s'", typ, format)
	}

	var err error
	var tags []string
	var validator *RegexValidator
	if validations != nil {
		if typ == "string" {
			tags, validator, err = makeStringValidations(name, schema.Format, validations)
			if err != nil {
				return nil, errors.Wrap(err, "error creating string restrictions")
			}
		} else if typ == "integer" {
			tags = generateIntegerRestriction(validations.Minimum, validations.Maximum, validations.ExclusiveMinimum, validations.ExclusiveMaximum)
		}
	}

	if schema.Type == "array" {
		typ := New(Simple, "", goType, true, false, tags...)
		typ.Validator = validator
		return NewArray(name, required, typ), nil
	} else {
		typ := New(Simple, name, goType, required, false, tags...)
		typ.Validator = validator
		return typ, nil
	}
}

func FromSchema(name string, schema *spec.Schema, required bool, findSchemaFunc func(name string) *spec.Schema) (*Type, error) {

	if schema == nil {
		return nil, errors.Errorf("'%s' has no schema", name)
	}

	if !schema.Ref.GetPointer().IsEmpty() {
		return referedType(name, schema, required, findSchemaFunc)
	}

	if len(schema.Enum) != 0 {

		return NewEnum(name, required, schema.Enum)

	} else if schema.Type.Contains("string") {

		validations := &spec.CommonValidations{
			ExclusiveMaximum: schema.ExclusiveMaximum,
			ExclusiveMinimum: schema.ExclusiveMinimum,
			MaxItems:         schema.MaxItems,
			MaxLength:        schema.MaxLength,
			Maximum:          schema.Maximum,
			MinItems:         schema.MinItems,
			MinLength:        schema.MinLength,
			Minimum:          schema.Minimum,
			MultipleOf:       schema.MultipleOf,
			Pattern:          schema.Pattern,
			UniqueItems:      schema.UniqueItems,
		}

		tags, validator, err := makeStringValidations(name, schema.Format, validations)
		if err != nil {
			return nil, errors.Wrap(err, "error creating string restrictions")
		}

		typ := New(Simple, name, StringType, required, false, tags...)
		if validator != nil {
			typ.Validator = validator
		}

		return typ, nil

	} else if schema.Type.Contains("file") {

		return NewFile(name, required), nil

	} else if schema.Type.Contains("boolean") {

		return New(Simple, name, BooleanType, required, false), nil

	} else if schema.Type.Contains("integer") {

		tags := generateIntegerRestriction(schema.Minimum, schema.Maximum, schema.ExclusiveMinimum, schema.ExclusiveMaximum)

		return New(Simple, name, ConvertSimpleType("integer", schema.Format), required, false, tags...), nil

	} else if schema.Type.Contains("number") {

		return New(Simple, name, Float64Type, required, false), nil

	} else if schema.Type.Contains("array") {

		elementType, err := FromSchema("", schema.Items.Schema, true, findSchemaFunc)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating array '%s'", name)
		}

		return NewArray(name, required, elementType), nil

	} else if len(schema.Type) == 0 || schema.Type.Contains("object") {

		if len(schema.Properties) == 0 && schema.AdditionalProperties != nil {

			elementType, err := FromSchema("", schema.AdditionalProperties.Schema, true, findSchemaFunc)
			if err != nil {
				return nil, errors.Wrapf(err, "error generating map '%s'", name)
			}

			return NewMap(name, required, elementType), nil

		} else {

			if name == "" {
				n := atomic.AddInt32(&objects, 1)
				name = "Object" + strconv.Itoa(int(n))
			}

			return objectType(name, schema, required, findSchemaFunc)
		}
	}

	return nil, errors.Errorf("unknown schema type '%s'", schema.Type[0])
}

func referedType(name string, schema *spec.Schema, required bool, findSchemaFunc func(name string) *spec.Schema) (*Type, error) {

	ref := schema.Ref.GetPointer().DecodedTokens()
	if len(ref) == 2 && ref[0] == "definitions" {

		referencedSchema := findSchemaFunc(ref[1])
		if referencedSchema == nil {
			return nil, errors.Errorf("schema ref '%s' not found", schema.Ref.GetURL())
		}

		if referencedSchema.Enum != nil && len(referencedSchema.Enum) != 0 {
			return New(Simple, name, identifier.MakeIdentifier(strings.Title(ref[1])), required, false), nil
		} else {
			return New(Simple, name, identifier.MakeIdentifier(strings.Title(ref[1])), required, true), nil
		}

	} else {
		return nil, errors.Errorf("unknown schema ref '%s'", schema.Ref.GetURL())
	}
}

var objects int32

func objectType(name string, schema *spec.Schema, required bool, findSchemaFunc func(name string) *spec.Schema) (*Type, error) {

	requiredProps := make(map[string]bool)
	for _, propRequired := range schema.Required {
		requiredProps[propRequired] = true
	}

	var propNames []string
	for propName := range schema.Properties {
		propNames = append(propNames, propName)
	}
	sort.Strings(propNames)

	typ := NewObject(name, required)
	for _, propName := range propNames {

		prop := schema.Properties[propName]
		isRequired := requiredProps[propName]

		if prop.Type.Contains("file") {
			return nil, errors.Errorf("error generating object '%s' (contains a file)", name)
		}

		elementType, err := FromSchema("", &prop, isRequired, findSchemaFunc)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating object '%s'", name)
		}

		typ.AddElement(strings.Title(identifier.MakeIdentifier(propName)), elementType, propName)
	}

	return typ, nil
}

func (typ *Type) AddElement(name string, t *Type, serialized string) {

	typ.Elements = append(typ.Elements, Element{Name: name, Type: t, Serialized: serialized})
}

func (typ *Type) GetValidators() []*RegexValidator {

	validators := []*RegexValidator{}

	if typ.Validator != nil {
		validators = append(validators, typ.Validator)
	}

	for _, element := range typ.Elements {
		if element.Type != nil {
			validators = append(validators, element.Type.GetValidators()...)
		}
	}

	return validators
}

func (typ *Type) WriteTo(file *file.File) error {

	if typ.Name == "" {
		return nil
	}

	if file.HasType(typ.Name) {
		log.Warn(fmt.Sprintf("suppressed duplicate of type '%s'", typ.Name))
		return nil
	}

	if typ.Composit == Simple || typ.Composit == File {

		file.Type().Id(typ.Name).Id(typ.Type)

	} else if typ.Composit == Object {

		var properties []jen.Code

		for _, element := range typ.Elements {

			if element.Type == nil {
				continue
			}

			err := element.Type.WriteTo(file)
			if err != nil {
				return errors.Wrapf(err, "error writing type of element %s", element.Name)
			}

			code := jen.Id(element.Name)

			if element.Type.Composit == Simple || element.Type.Composit == Enum {
				if !element.Type.Required {
					code.Op("*")
				}
			} else if element.Type.Composit == Array {
				code.Index()
			} else if element.Type.Composit == Map {
				code.Map(jen.String())
			}

			if element.Type.Name != "" {
				code.Id(element.Type.Name)
			} else if element.Type.Type != "" {
				code.Id(element.Type.Type)
			}
			if element.Type.Validation != "" || element.Serialized != "" {

				tags := map[string]string{}
				if element.Type.Validation != "" {
					tags["validate"] = element.Type.Validation
				}

				if element.Serialized != "" {
					serialized := element.Serialized
					if element.Type.Required {
						serialized += ",required"
					} else {
						serialized += ",omitempty"
					}

					tags["xml"] = serialized
					tags["json"] = serialized
					tags["bson"] = serialized
				}

				code.Tag(tags)
			}

			properties = append(properties, code)
		}

		file.AddType(typ.Name).Struct(properties...).Line()

	} else if typ.Composit == Array {

		if len(typ.Elements) > 0 && typ.Elements[0].Type != nil {
			err := typ.Elements[0].Type.WriteTo(file)
			if err != nil {
				return errors.Wrapf(err, "error writing type '%s'", typ.Name)
			}
			if typ.Elements[0].Type.Composit == Enum {
				return nil
			}
		}

		file.AddType(typ.Name).Index().Id(typ.Type)

	} else if typ.Composit == Map {

		if len(typ.Elements) > 0 && typ.Elements[0].Type != nil {
			err := typ.Elements[0].Type.WriteTo(file)
			if err != nil {
				return errors.Wrapf(err, "error writing type '%s'", typ.Name)
			}
		}

		file.AddType(typ.Name).Map(jen.String()).Id(typ.Type)

	} else if typ.Composit == Enum {

		file.AddType(typ.Name).Id(typ.Type)

		var constants []jen.Code
		for _, element := range typ.Elements {
			if element.Type != nil {
				constants = append(constants, jen.Id(element.Name).Id(element.Type.Type).Op("=").Id(element.Type.Name))
			}
		}

		file.Const().Defs(constants...)
	}

	return nil
}
