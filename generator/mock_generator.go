package generator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GoClientMockGenerator struct {
	*GoGenerator
}

func NewGoClientMockGenerator(spec *openapi.Spec) *GoClientMockGenerator {

	return &GoClientMockGenerator{
		GoGenerator: NewGoGenerator(spec),
	}
}

func (gen *GoClientMockGenerator) Generate(path, pckg string) error {

	file := jen.NewFile(pckg)
	nameOfClient := identifier.MakeIdentifier(strings.Title(gen.Spec.Info().Title + "ClientMock"))
	gen.generateConstructor(nameOfClient, file)

	var attributes []jen.Code
	if err := gen.WalkOperations(func(operation *Operation) error {
		err := gen.generateOperation(nameOfClient, operation, file)
		if err != nil {
			return errors.Wrapf(err, "error generating code for %s '%v'", operation.Method, operation.Operation)
		}

		attributes = append(attributes, jen.Id(strings.Title(operation.ID)+"StatusCode").Int())
		return nil
	}); err != nil {
		return errors.Wrap(err, "error generating operations")
	}

	file.Type().Id(nameOfClient).Struct(attributes...).Line()

	err := file.Save(path)
	if err != nil {
		return errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}
	return nil
}

func (gen *GoClientMockGenerator) generateConstructor(name string, file *jen.File) {

	file.Func().Id("New"+name).Params(
		jen.Id("httpClient").Op("*").Qual("net/http", "Client"),
		jen.Id("baseUrl").String(),
		jen.Id("ctx").Op("...").Qual("context", "Context"),
	).Params(
		jen.Op("*").Id(name),
	).Block(
		jen.Return(jen.Op("&").Id(name).Values()),
	).Line()
}

func (gen *GoClientMockGenerator) generateOperation(structName string, operation *Operation, file *jen.File) error {

	if operation.Description != "" {
		file.Comment(operation.Description)
	} else {
		file.Line()
	}

	var code []jen.Code
	walkResponses(operation.Responses.StatusCodeResponses, func(status int, response spec.Response) {
		handler, err := gen.makeStatusCode("client", operation, status, response)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("error generating code for %s '%v' status code %d", operation.Method, operation.Operation, status))
		} else {
			code = append(code, handler)
		}
	})
	code = append(code, jen.Return(jen.Nil(), jen.Nil()))

	name := strings.Title(operation.ID)
	file.Func().Params(
		jen.Id("client").Op("*").Id(structName),
	).Id(name).Params(
		jen.Id("request").Op("*").Id(name+"Request"),
	).Params(
		jen.Id(name+"Response"), jen.Id("error"),
	).Block(
		code...,
	).Line()

	return nil
}

func (gen *GoClientMockGenerator) makeStatusCode(client string, operation *Operation, status int, response spec.Response) (jen.Code, error) {

	responseStruct := strings.Title(operation.ID) + strconv.Itoa(status) + "Response"
	code := []jen.Code{
		jen.Id("response").Op(":=").New(jen.Id(responseStruct)),
		jen.Return(jen.Id("response"), jen.Nil()),
	}

	if response.Schema != nil && operation.HasProduce(ContentTypeApplicationJson) {
		data, ok := response.ResponseProps.Examples[ContentTypeApplicationJson]
		if ok {
			raw, ok := data.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("unsupported example data format (%v)", data)
			}

			jsonData, err := json.Marshal(raw)
			if err != nil {
				return nil, errors.Errorf("could not marshal data (%v)", data)
			}

			code = []jen.Code{
				jen.Id("data").Op(":=").Lit(string(jsonData)),
				jen.Id("response").Op(":=").New(jen.Id(responseStruct)),
				jen.Id("responseBody").Op(":=").Op("&").Id("response").Dot("Body"),
				jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(jen.Index().Byte().Parens(jen.Id("data")), jen.Id("responseBody")),
				jen.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Id("err"))),
				jen.Return(jen.Id("response"), jen.Nil()),
			}
		} else {
			log.Info(fmt.Sprintf("json example blob is missing for operation %s and status code %d", strings.Title(operation.ID), status))
		}
	}

	handler := jen.If(jen.Id("client")).Dot(strings.Title(operation.ID) + "StatusCode").Op("==").Lit(status).Block(code...)
	return handler, nil
}
