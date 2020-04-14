package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ExperienceOne/apikit/generator/file"
	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/ExperienceOne/apikit/generator/types"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ValidatorMap map[string]string

func NewValidatorMap() ValidatorMap {

	return make(ValidatorMap)
}

func (m ValidatorMap) Add(v *types.RegexValidator) {

	m[v.Tag] = v.Regex
}

func (m ValidatorMap) AddAll(v []*types.RegexValidator) {

	for _, validator := range v {
		m.Add(validator)
	}
}

type goTypesGenerator struct {
	*GoGenerator
	validators ValidatorMap
}

func NewGoTypesGenerator(spec *openapi.Spec) *goTypesGenerator {
	return &goTypesGenerator{
		GoGenerator: NewGoGenerator(spec),
	}
}

func (gen *goTypesGenerator) Generate(path, pckg string) (ValidatorMap, error) {

	gen.validators = NewValidatorMap()

	file := file.NewFile(pckg)
	file.Var().Id("contentTypesForFiles").Op("=").Index().String().ValuesFunc(func(group *jen.Group) {
		for _, contentType := range ContentTypesForFiles {
			group.Add(jen.Lit(contentType))
		}
	})

	if err := gen.walkDefinitions(func(name string, definition *spec.Schema) (bool, error) {
		err := gen.generateDefinition(name, definition, file)
		if err != nil {
			return false, errors.Wrapf(err, "error generating definition '%s'", name)
		}
		return false, nil
	}); err != nil {
		return nil, errors.Wrap(err, "error generating definitions")
	}

	if err := gen.walkParameters(func(name string, parameter *spec.Parameter) (bool, error) {
		err := gen.generateParameter(name, parameter, file)
		if err != nil {
			return false, errors.Wrapf(err, "error generating parameter '%s'", name)
		}
		return false, nil
	}); err != nil {
		return nil, errors.Wrap(err, "error generating parameters")
	}

	if err := gen.WalkOperations(func(operation *Operation) error {

		var parameters []spec.Parameter
		parameters = append(parameters, operation.Path.Parameters...)
		parameters = append(parameters, operation.Parameters...)

		parametersBucket := NewParameterBucket(operation.HasConsume(ContentTypeApplicationFormUrlencoded))
		bucket, err := gen.PopulateParametersBucket(parametersBucket, parameters, operation.Security)
		if err != nil {
			return err
		}

		if err := gen.generateRequest(strings.Title(operation.ID+"Request"), bucket, file); err != nil {
			return errors.Wrapf(err, "error generating request for '%s'", operation.ID)
		}

		gen.generateResponses(operation, file)

		return nil
	}); err != nil {
		return nil, err
	}

	err := file.Save(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}

	return gen.validators, nil
}

func (gen *goTypesGenerator) walkParameters(handler func(name string, parameter *spec.Parameter) (bool, error)) error {

	parameters := gen.Spec.Parameters()

	var names []string
	for name := range parameters {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		parameter := parameters[name]
		abort, err := handler(name, &parameter)
		if err != nil {
			return err
		}

		if abort {
			return nil
		}
	}

	return nil
}

func (gen *goTypesGenerator) generateParameter(parameterName string, parameter *spec.Parameter, file *file.File) error {

	typ, err := gen.makeParameter(parameter)
	if err != nil {
		return errors.Wrapf(err, "error generating parameter '%s'", parameterName)
	}

	for _, validator := range typ.GetValidators() {
		gen.validators[validator.Tag] = validator.Regex
	}

	gen.validators.AddAll(typ.GetValidators())

	if parameter.Description != "" {
		file.Comment(parameter.Description)
	}

	err = typ.WriteTo(file)
	if err != nil {
		return errors.Wrapf(err, "error writing parameter '%s'", parameterName)
	}

	return nil
}

func (gen *goTypesGenerator) walkDefinitions(handler func(name string, definition *spec.Schema) (bool, error)) error {

	definitions := gen.Spec.Definitions()

	i := 0
	names := make([]string, len(definitions))
	for name := range definitions {
		names[i] = name
		i++
	}
	sort.Strings(names)

	for _, name := range names {
		definition := definitions[name]
		abort, err := handler(name, &definition)
		if err != nil {
			return err
		}

		if abort {
			return nil
		}
	}

	return nil
}

func (gen *goTypesGenerator) generateDefinition(definitionName string, definition *spec.Schema, file *file.File) error {

	typ, err := types.FromSchema(strings.Title(identifier.MakeIdentifier(definitionName)), definition, true, gen.FindDefinitionByName)
	if err != nil {
		return errors.Wrapf(err, "error generating definition '%s'", definitionName)
	}

	gen.validators.AddAll(typ.GetValidators())

	if definition.Description != "" {
		file.Comment(definition.Description)
	}

	err = typ.WriteTo(file)
	if err != nil {
		return errors.Wrapf(err, "error writing definition '%s'", definitionName)
	}

	return nil
}

func (gen *goTypesGenerator) FindDefinitionByName(name string) *spec.Schema {

	var definition *spec.Schema

	if err := gen.walkDefinitions(func(n string, d *spec.Schema) (bool, error) {
		if name == n {
			definition = d
			return true, nil
		}
		return false, nil
	}); err != nil {
		log.WithError(err).Error(fmt.Sprintf("error searching for definition '%s'", name))
		return nil
	}

	return definition
}

func (gen *goTypesGenerator) generateRequest(requestName string, bucket *ParametersBucket, file *file.File) error {

	if bucket.HasBody && bucket.HasFormData {
		return errors.New("definition of both form data and body type is not allowed in a request")
	}

	request := types.NewObject(requestName, true)

	if bucket.HasFormData {
		formData := types.NewObject(requestName+"FormData", true)

		for _, file := range bucket.FormDataFiles {
			formData.AddElement(identifier.MakeIdentifier(strings.Title(file.Name)), types.New(types.Simple, "", "MimeFile", file.Required, false), "")
		}

		for _, field := range bucket.FormData {
			if typ, err := types.FromSimpleSchema("", &field.SimpleSchema, field.Required, nil); err != nil {
				return errors.Wrap(err, "error generating form data")
			} else {
				gen.validators.AddAll(typ.GetValidators())
				formData.AddElement(identifier.MakeIdentifier(strings.Title(field.Name)), typ, "")
			}
		}

		request.AddElement("FormData", formData, "")
	}

	var parameters []*spec.Parameter
	parameters = append(parameters, bucket.Path...)
	parameters = append(parameters, bucket.Query...)
	parameters = append(parameters, bucket.Header...)
	if bucket.HasURLEncoded {
		parameters = append(parameters, bucket.FormData...)
	}
	parameters = append(parameters, bucket.Body...)

	for _, parameter := range parameters {
		typ, err := gen.makeParameter(parameter)
		if err != nil {
			return err
		}
		request.AddElement(identifier.MakeIdentifier(strings.Title(parameter.Name)), typ, "")
	}

	for _, parameter := range bucket.Security {
		parameterName := identifier.MakeIdentifier(strings.Title(parameter.Name))

		if parameter.Type == openapi.ApiKey.String() && parameter.In == openapi.Header.String() {
			request.AddElement(parameterName, types.New(types.Simple, "", types.StringType, true, false), "")
		} else if parameter.Type == openapi.Basic.String() {
			if parameterName == "" {
				parameterName = openapi.BasicAuthHeaderName
			}
			request.AddElement(parameterName, types.New(types.Simple, "", types.StringType, true, false), "")
		} else {
			return errors.Errorf("unsupported security scheme type '%s' in '%s'", parameter.Type, parameter.In)
		}
	}

	err := request.WriteTo(file)
	if err != nil {
		return errors.Wrapf(err, "error writing request '%s'", requestName)
	}
	return nil
}

func (gen *goTypesGenerator) makeParameter(parameter *spec.Parameter) (*types.Type, error) {

	if parameter.Schema != nil {

		typ, err := types.FromSchema("", parameter.Schema, parameter.Required, gen.FindDefinitionByName)
		if err != nil {
			return nil, err
		}
		gen.validators.AddAll(typ.GetValidators())
		return typ, nil

	} else if len(parameter.Enum) > 0 {

		enumType, err := types.NewEnum(strings.Title(parameter.Name), parameter.Required, parameter.Enum)
		if err != nil {
			return nil, err
		}

		if parameter.Type == "array" {
			return types.NewArray(strings.Title(parameter.Name), parameter.Required, enumType), nil
		} else {
			return enumType, nil
		}

	} else {
		typ, err := types.FromSimpleSchema("", &parameter.SimpleSchema, parameter.Required, &parameter.CommonValidations)
		gen.validators.AddAll(typ.GetValidators())
		return typ, err
	}
}

func (gen *goTypesGenerator) generateResponses(operation *Operation, file *file.File) {

	interfaceName := strings.Title(operation.ID + "Response")
	file.Type().Id(interfaceName).Interface(
		jen.Id("is"+interfaceName).Params(),
		jen.Id("StatusCode").Params().Int(),
		jen.Id("write").Params(jen.Id("response").Qual("net/http", "ResponseWriter")).Error(),
	).Line()

	walkResponses(operation.Responses.StatusCodeResponses, func(statusCode int, response spec.Response) {

		responseName := strings.Title(fmt.Sprintf("%s%dResponse", operation.ID, statusCode))
		responseType, err := gen.makeResponse(responseName, response.Schema, response.Headers)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("error generating %d response for '%s'", statusCode, operation.ID))
			return
		}

		if response.Description != "" {
			file.Comment(response.Description)
		}
		err = responseType.WriteTo(file)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("error writing %d response for '%s'", statusCode, operation.ID))
			return
		}

		file.Func().Params(jen.Id("r").Op("*").Id(responseName)).Id("is" + interfaceName).Params().Block().Line()

		file.Func().Params(jen.Id("r").Op("*").Id(responseName)).Id("StatusCode").Params().Int().Block(
			jen.Return(jen.Lit(statusCode)),
		).Line()

		var (
			writeFunc []jen.Code
			headers   []string
		)

		for header := range response.Headers {
			headers = append(headers, header)
		}
		sort.Strings(headers)

		for _, header := range headers {
			field := identifier.MakeIdentifier(strings.Title(header))
			// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
			writeHeader := jen.Id("response").Dot("Header").Call().Index(jen.Lit(header)).Op("=").Index().String().Values(jen.Id("toString").Call(jen.Id("r").Dot(field)))
			writeFunc = append(writeFunc, writeHeader)
		}

		// Behavior of content type
		// Type
		// 	https://tools.ietf.org/search/rfc2616#section-7.2
		// Entity body
		// 	https://tools.ietf.org/search/rfc2616#section-7.2.1

		returnInternalServerError := jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError")))

		if response.Schema != nil && response.Schema.Type.Contains("file") && operation.HasProduces(ContentTypesForFiles...) {

			writeFile := jen.If(jen.List(jen.Id("_"), jen.Id("err")).Op(":=").Qual("io", "Copy").Call(jen.Id("response"), jen.Id("r").Dot("Body")), jen.Id("err").Op("!=").Nil()).Block(
				returnInternalServerError,
			)
			writeFunc = append(writeFunc, writeFile)

			closeBody := jen.If(jen.Id("err").Op(":=").Id("r").Dot("Body").Dot("Close").Call(), jen.Id("err").Op("!=").Nil()).Block(
				returnInternalServerError,
			)
			writeFunc = append(writeFunc, closeBody)

		} else if response.Schema != nil && operation.HasProduce(ContentTypeApplicationJson) {

			writeJson := jen.If(jen.Id("err").Op(":=").Id("serveJson").Call(jen.Id("response"), jen.Lit(statusCode), jen.Id("r").Dot("Body")), jen.Id("err").Op("!=").Nil()).Block(
				returnInternalServerError,
			)
			writeFunc = append(writeFunc, writeJson)

		} else if response.Schema != nil && operation.HasProduce(ContentTypeApplicationHalJson) {

			writeJson := jen.If(jen.Id("err").Op(":=").Id("serveHalJson").Call(jen.Id("response"), jen.Lit(statusCode), jen.Id("r").Dot("Body")), jen.Id("err").Op("!=").Nil()).Block(
				returnInternalServerError,
			)
			writeFunc = append(writeFunc, writeJson)

		} else {
			// don't set a content type for an nil schema
			clearContentType := jen.Id("response").Dot("Header").Call().Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values()
			writeFunc = append(writeFunc, clearContentType)

			writeStatusCode := jen.Id("response").Dot("WriteHeader").Call(jen.Lit(statusCode))
			writeFunc = append(writeFunc, writeStatusCode)
		}

		writeFunc = append(writeFunc, jen.Return(jen.Nil()))

		file.Func().Params(jen.Id("r").Op("*").Id(responseName)).Id("write").Params(jen.Id("response").Qual("net/http", "ResponseWriter")).Error().Block(
			writeFunc...,
		).Line()
	})
}

func (gen *goTypesGenerator) makeResponse(responseName string, body *spec.Schema, headers map[string]spec.Header) (*types.Type, error) {

	response := types.NewObject(responseName, true)

	if body != nil {
		bodyType, err := types.FromSchema("", body, true, gen.FindDefinitionByName)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating response '%s'", responseName)
		}

		gen.validators.AddAll(bodyType.GetValidators())
		response.AddElement("Body", bodyType, "")
	}

	var headerNames []string
	for headerName := range headers {
		headerNames = append(headerNames, headerName)
	}
	sort.Strings(headerNames)

	for _, name := range headerNames {
		header := headers[name]
		// OpenAPI spec does not support not-required response headers
		typ, err := types.FromSimpleSchema("", &header.SimpleSchema, true, &header.CommonValidations)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating response '%s'", responseName)
		}
		gen.validators.AddAll(typ.GetValidators())
		response.AddElement(identifier.MakeIdentifier(strings.Title(name)), typ, "")
	}

	return response, nil
}
