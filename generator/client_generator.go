package generator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ExperienceOne/apikit/generator/xhttp"

	"github.com/ExperienceOne/apikit/generator/file"
	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/ExperienceOne/apikit/generator/stringutil"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type goClientGenerator struct {
	*GoGenerator
	mockGoClientGenerator *mockGoClientGenerator
}

func NewGoClientGenerator(spec *openapi.Spec) *goClientGenerator {

	return &goClientGenerator{
		GoGenerator:           NewGoGenerator(spec),
		mockGoClientGenerator: NewMockGoClientGenerator(),
	}
}

func (gen goClientGenerator) clientName() string {
	return stringutil.UnTitle(identifier.MakeIdentifier(gen.Spec.Info().Title + "Client"))
}

func (gen *goClientGenerator) Generate(path, pckg string, generatePrometheus bool, generateMocks bool) error {

	file := file.NewFile(pckg)

	var clientMethods []jen.Code
	operations := make([]*Operation, 0)
	if err := gen.WalkOperations(func(operation *Operation) error {
		methodDef := jen.Id(strings.Title(operation.ID)).Params(jen.Id("request").Op("*").Id(strings.Title(operation.ID+"Request"))).Params(jen.Id(strings.Title(operation.ID+"Response")), jen.Error())
		clientMethods = append(clientMethods, methodDef)
		operations = append(operations, operation)
		return nil
	}); err != nil {
		return errors.Wrap(err, "error generating operations")
	}

	file.Type().Id(strings.Title(gen.clientName())).Interface(clientMethods...)

	gen.generateConstructor(gen.clientName(), pckg, operations, generatePrometheus, file)

	clientMembers := []jen.Code{
		jen.Id("baseURL").String(),
		jen.Id("hooks").Id("HooksClient"),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("httpClient").Op("*").Id("httpClientWrapper"),
		jen.Id("xmlMatcher").Op("*").Qual("regexp", "Regexp"),
	}

	if generatePrometheus {
		clientMembers = append(clientMembers, jen.Id("prometheusHandler").Op("*").Id("PrometheusHandler"))
	}

	file.Type().Id(gen.clientName()).Struct(clientMembers...)

	if err := gen.WalkOperations(func(operation *Operation) error {

		err := gen.generateOperation(operation, gen.clientName(), generatePrometheus, file)
		if err != nil {
			return errors.Wrapf(err, "error generating code for %s '%v'", operation.Method, operation.Operation)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "error generating operations")
	}

	err := file.Save(path)
	if err != nil {
		return errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}

	if generateMocks {
		lastSlash := -1
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' || path[i] == '\\' {
				lastSlash = i
				break
			}
		}
		err := gen.mockGoClientGenerator.Generate(gen.clientName(), path[:lastSlash])
		if err != nil {
			return errors.Wrap(err, "error generating mock client")
		}
	}

	return nil
}

func (gen *goClientGenerator) generateConstructor(nameOfClient, nameOfPackage string, operations []*Operation, generatePrometheus bool, file *file.File) {

	interfaceName := strings.Title(nameOfClient)

	file.Func().Id("New"+interfaceName).Params(jen.Id("httpClient").Op("*").Qual("net/http", "Client"), jen.Id("baseUrl").String(), jen.Id("options").Id("Opts")).Id(interfaceName).BlockFunc(func(stmts *jen.Group) {

		if generatePrometheus {
			stmts.Id("namespace").Op(":=").Lit(nameOfPackage)
			stmts.Id("h").Op(":=").Id("NewPrometheusHandler").Call(jen.Op("&").Id("namespace")).Line()

			for _, op := range operations {
				route := strings.Replace(strings.Replace(op.Route, "{", "<", -1), "}", ">", -1)

				stmts.Id("h").Dot("InitMetric").Call(
					jen.Lit(route),
					jen.Lit(op.Method),
				)
			}
			stmts.Line()
		}

		handlerParameters := []jen.Code{
			jen.Id("httpClient").Op(":").Id("newHttpClientWrapper").Call(jen.Id("httpClient"), jen.Id("baseUrl")),
			jen.Id("baseURL").Op(":").Id("baseUrl"),
			jen.Id("hooks").Op(":").Id("options").Dot("Hooks"),
			jen.Id("ctx").Op(":").Id("options").Dot("Ctx"),
			jen.Id("xmlMatcher").Op(":").Qual("regexp", "MustCompile").Call(jen.Lit(ContentTypeApplicationXMLPattern)),
		}

		if generatePrometheus {
			handlerParameters = append(handlerParameters, jen.Id("prometheusHandler").Op(":").Id("h"))
		}

		stmts.Return(jen.Op("&").Id(nameOfClient).Values(handlerParameters...))
	}).Line()
}

func (gen *goClientGenerator) generateOperation(operation *Operation, nameOfClient string, generatePrometheus bool, file *file.File) error {

	if operation.Description != "" {
		file.Comment(operation.Description)
	}

	var parameters []spec.Parameter
	parameters = append(parameters, operation.Path.Parameters...)
	parameters = append(parameters, operation.Parameters...)

	bucket, err := gen.PopulateParametersBucket(NewParameterBucket(operation.HasConsume(ContentTypeApplicationFormUrlencoded)), parameters, operation.Security)
	if err != nil {
		return err
	}

	file.Func().Params(jen.Id("client").Op("*").Id(nameOfClient)).Id(strings.Title(operation.ID)).Params(jen.Id("request").Op("*").Id(strings.Title(operation.ID+"Request"))).Params(jen.Id(strings.Title(operation.ID+"Response")), jen.Error()).BlockFunc(func(stmts *jen.Group) {

		if !operation.HasValidConsumes() || !operation.HasValidProduces() {
			gen.generateNotSupported(operation, "no supported content type", true, stmts)
			return
		}

		if match, valid := hasFileEndpointValidProduce(operation); match && valid == 0 {
			gen.generateNotSupported(operation, strings.Join(operation.Consumes, ","), true, stmts)
			return
		}

		if len(bucket.Body) != 0 && operation.HasConsume(ContentTypeApplicationFormUrlencoded) {
			gen.generateNotSupported(operation, "invalid definition (both body type & application/x-www-form-urlencoded)", true, stmts)
			return
		}

		stmts.If(jen.Id("request").Op("==").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("newRequestObjectIsNilError")),
		)

		stmts.Id("path").Op(":=").Lit(operation.Route)
		stmts.Id("method").Op(":=").Lit(operation.Method)
		stmts.Id("endpoint").Op(":=").Id("client").Dot("baseURL").Op("+").Id("path")
		stmts.Id("httpContext").Op(":=").Id("newHttpContextWrapper").Call(jen.Id("client").Dot("ctx"))

		for _, param := range bucket.Path {
			stmts.Id("endpoint").Op("=").Qual("strings", "Replace").Call(
				jen.Id("endpoint"),
				jen.Lit("{"+param.Name+"}"),
				jen.Qual("net/url", "QueryEscape").Call(jen.Id("toString").Call(jen.Id("request").Dot(identifier.MakeIdentifier(strings.Title(param.Name))))),
				jen.Lit(1),
			)
		}

		queryParamCount := 0
		if len(bucket.Query) != 0 {
			queryParamCount = gen.generateQueryString(bucket.Query, "query", queryParamCount, stmts)
			stmts.Id("encodedQuery").Op(":=").Id("query").Dot("Encode").Call()
			stmts.If(jen.Id("encodedQuery").Op("!=").Lit("")).Block(
				jen.Id("endpoint").Op("+=").Lit("?").Op("+").Id("encodedQuery"),
			)
		}

		if operation.HasConsume(ContentTypeApplicationFormUrlencoded) {

			gen.generateQueryString(bucket.FormData, "queryInBody", queryParamCount, stmts)
			stmts.Id("encodedQueryInBody").Op(":=").Id("queryInBody").Dot("Encode").Call()
			stmts.Id("formData").Op(":=").Qual("bytes", "NewBufferString").Call(jen.Id("encodedQueryInBody"))
			stmts.List(jen.Id("httpRequest"), jen.Id("reqErr")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("method"), jen.Id("endpoint"), jen.Id("formData"))

		} else if operation.HasConsume(ContentTypeMultipartFormData) {

			stmts.Id("formData").Op(":=").New(jen.Qual("bytes", "Buffer"))
			stmts.Id("bodyWriter").Op(":=").Qual("mime/multipart", "NewWriter").Call(jen.Id("formData"))

			for i, param := range bucket.FormDataFiles {

				fileField := strings.Title(identifier.MakeIdentifier(param.Name))

				writerErr := "writerErr" + strconv.Itoa(i)
				fileWriter := "fileWriter" + strconv.Itoa(i)

				createFormFile := jen.List(jen.Id(fileWriter), jen.Id(writerErr)).Op(":=").Id("bodyWriter").Dot("CreateFormFile").Call(jen.Lit(param.Name), jen.Lit(param.Name))
				createFormFileError := jen.If(jen.Id(writerErr).Op("!=").Nil()).Block(
					jen.Id("bodyWriter").Dot("Close").Call(),
					jen.Return(jen.Nil(), jen.Id(writerErr)),
				)

				copyErr := "copyFileErr" + strconv.Itoa(i)
				copyFormFile := jen.List(jen.Id("_"), jen.Id(copyErr)).Op(":=").Qual("io", "Copy").Call(jen.Id(fileWriter), jen.Id("request").Dot("FormData").Dot(fileField).Dot("Content"))
				copyFormFileError := jen.If(jen.Id(copyErr).Op("!=").Nil()).Block(
					jen.Id("bodyWriter").Dot("Close").Call(),
					jen.Return(jen.Nil(), jen.Id(copyErr)),
				)

				if param.Required {
					stmts.Block(createFormFile, createFormFileError, copyFormFile, copyFormFileError)
				} else {
					stmts.If(jen.Id("request").Dot("FormData").Dot(fileField).Op("!=").Nil()).Block(
						createFormFile, createFormFileError, copyFormFile, copyFormFileError,
					)
				}
			}

			for i, param := range bucket.FormData {

				fieldValue := "fieldData" + strconv.Itoa(i)
				dataField := strings.Title(identifier.MakeIdentifier(param.Name))

				toString := jen.Id(fieldValue).Op(":=").Id("toString").Call(jen.Id("request").Dot("FormData").Dot(dataField))

				fieldErr := "fieldErr" + strconv.Itoa(i)
				fieldWriter := "fieldWriter" + strconv.Itoa(i)

				createFormField := jen.List(jen.Id(fieldWriter), jen.Id(fieldErr)).Op(":=").Id("bodyWriter").Dot("CreateFormField").Call(jen.Lit(param.Name))
				createFormFieldError := jen.If(jen.Id(fieldErr).Op("!=").Nil()).Block(
					jen.Id("bodyWriter").Dot("Close").Call(),
					jen.Return(jen.Nil(), jen.Id(fieldErr)),
				)

				writeErr := "writeFieldErr" + strconv.Itoa(i)
				writeFormField := jen.List(jen.Id("_"), jen.Id(writeErr)).Op(":=").Id(fieldWriter).Dot("Write").Call(jen.Index().Byte().Params(jen.Id(fieldValue)))
				writeFormFieldError := jen.If(jen.Id(writeErr).Op("!=").Nil()).Block(
					jen.Id("bodyWriter").Dot("Close").Call(),
					jen.Return(jen.Nil(), jen.Id(writeErr)),
				)

				if param.Required {
					stmts.Add(toString, createFormField, createFormFieldError, writeFormField, writeFormFieldError)
				} else {
					stmts.If(jen.Id("request").Dot("FormData").Dot(dataField)).Op("!=").Nil().Block(
						toString, createFormField, createFormFieldError, writeFormField, writeFormFieldError,
					)
				}
			}

			stmts.Id("contentType").Op(":=").Id("bodyWriter").Dot("FormDataContentType").Call()
			stmts.Id("bodyWriter").Dot("Close").Call()
			stmts.List(jen.Id("httpRequest"), jen.Id("reqErr")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("method"), jen.Id("endpoint"), jen.Id("formData"))

		} else if bucket.HasBody {

			var bodyName string
			for _, b := range bucket.Body {
				if b.In == openapi.Body.String() {
					bodyName = strings.Title(identifier.MakeIdentifier(b.Name))
					break
				}
			}

			if bodyName == "" {
				bodyName = "Body"
			}

			stmts.Id("jsonData").Op(":=").New(jen.Qual("bytes", "Buffer"))
			stmts.Id("encodeErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Id("jsonData")).Dot("Encode").Call(jen.Op("&").Id("request").Dot(bodyName))
			stmts.If(jen.Id("encodeErr").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("encodeErr")),
			)
			stmts.List(jen.Id("httpRequest"), jen.Id("reqErr")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("method"), jen.Id("endpoint"), jen.Id("jsonData"))

		} else {
			stmts.List(jen.Id("httpRequest"), jen.Id("reqErr")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Id("method"), jen.Id("endpoint"), jen.Nil())
		}

		stmts.If(jen.Id("reqErr").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("reqErr")),
		)

		// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
		if operation.HasConsume(ContentTypeApplicationFormUrlencoded) {
			stmts.Id("httpRequest").Dot("Header").Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values(jen.Id("contentTypeApplicationFormUrlencoded"))
		} else if operation.HasConsume(ContentTypeMultipartFormData) {
			stmts.Id("httpRequest").Dot("Header").Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values(jen.Id("contentType"))
		} else if operation.HasConsume(ContentTypeApplicationHalJson) {
			stmts.Id("httpRequest").Dot("Header").Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values(jen.Id("ContentTypeApplicationHalJson"))
		} else if bucket.HasBody {
			stmts.Id("httpRequest").Dot("Header").Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values(jen.Id("contentTypeApplicationJson"))
		}

		for _, param := range bucket.Security {

			if param.Type == openapi.ApiKey.String() && param.In == openapi.Header.String() {
				// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
				stmts.Id("httpRequest").Dot("Header").Index(jen.Lit(param.Name)).Op("=").Index().String().Values(jen.Id("request").Dot(identifier.MakeIdentifier(strings.Title(param.Name))))
			} else if param.Type == openapi.Basic.String() {
				paramName := param.Name
				if paramName == "" {
					paramName = openapi.BasicAuthHeaderName
				}
				// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
				stmts.Id("httpRequest").Dot("Header").Index(jen.Lit(paramName)).Op("=").Index().String().Values(jen.Id("request").Dot(identifier.MakeIdentifier(strings.Title(paramName))))
			} else {
				log.Error(fmt.Sprintf("unsupported security scheme type '%s' in '%s'", param.Type, param.In))
			}
		}

		for _, param := range bucket.Header {

			paramField := identifier.MakeIdentifier(strings.Title(param.Name))
			// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
			setHeader := jen.Id("httpRequest").Dot("Header").Index(jen.Lit(param.Name)).Op("=").Index().String().Values(jen.Id("toString").Call(jen.Id("request").Dot(paramField)))

			if param.Required {
				stmts.Add(setHeader)
			} else {
				stmts.If(jen.Id("request").Dot(paramField).Op("!=").Nil()).Block(setHeader)
			}
		}

		stmts.Comment("set all headers from client context")
		stmts.Id("err").Op(":=").Id("setRequestHeadersFromContext").Call(jen.Id("httpContext"), jen.Id("httpRequest").Dot("Header"))
		stmts.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		)

		if operation.HasValidProduces() {
			stmts.If(jen.Len(jen.Id("httpRequest").Dot("Header").Index(jen.Lit("accept"))).Op("==").Lit(0).Op("&&").Len(jen.Id("httpRequest").Dot("Header").Index(jen.Lit("Accept"))).Op("==").Lit(0)).Block(
				jen.Id("httpRequest").Dot("Header").Index(jen.Lit("Accept")).Op("=").Index().String().Values(jen.Lit(strings.Join(operation.Produces, ", "))),
			)
		}

		if generatePrometheus {
			stmts.Id("start").Op(":=").Qual("time", "Now").Call()
		}

		stmts.List(jen.Id("httpResponse"), jen.Id("err")).Op(":=").Id("client").Dot("httpClient").Dot("Do").Call(jen.Id("httpRequest"))
		stmts.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		)

		lastFileResponseIndex := -1
		walkResponses(operation, func(statusCode int, response spec.Response) {
			if response.Schema != nil && response.Schema.Type.Contains("file") && operation.HasProduces(ContentTypesForFiles...) {
				lastFileResponseIndex++
			}
		})

		// if the list of responses has a file response then don't close the http body response type here
		if lastFileResponseIndex == -1 {
			stmts.Defer().Id("httpResponse").Dot("Body").Dot("Close").Call()
		}

		if generatePrometheus {
			stmts.Id("client").Dot("prometheusHandler").Dot("HandleRequest").Call(
				jen.Id("path"),
				jen.Id("method"),
				jen.Id("httpResponse").Dot("StatusCode"),
				jen.Qual("time", "Since").Call(jen.Id("start")),
			)
		}

		var currentIndex int
		var addedDeferStatement bool
		walkResponses(operation, func(statusCode int, response spec.Response) {
			gen.generateResponse(operation, statusCode, response, stmts)
			// generate defer http.Response.Close statement after last file response handler
			if lastFileResponseIndex > -1 && currentIndex == lastFileResponseIndex && addedDeferStatement == false {
				addedDeferStatement = true
				stmts.Defer().Id("httpResponse").Dot("Body").Dot("Close").Call()
			}
			currentIndex++
		})

		stmts.If(jen.Id("client").Dot("hooks").Dot("OnUnknownResponseCode").Op("!=").Nil()).Block(
			jen.Id("message").Op(":=").Id("client").Dot("hooks").Dot("OnUnknownResponseCode").Call(jen.Id("httpResponse"), jen.Id("httpRequest")),
			jen.Return(jen.Nil(), jen.Id("newErrOnUnknownResponseCode").Call(jen.Id("message"))),
		)

		stmts.Return(jen.Nil(), jen.Id("newErrUnknownResponse").Call(jen.Id("httpResponse").Dot("StatusCode")))
	}).Line()
	return nil
}

func (gen *goClientGenerator) generateQueryString(params []*spec.Parameter, id string, lastIndex int, stmts *jen.Group) int {

	stmts.Id(id).Op(":=").Make(jen.Qual("net/url", "Values"))

	for _, param := range params {

		name := strings.Title(identifier.MakeIdentifier(param.Name))
		addParam := jen.Id(id).Dot("Add").Call(jen.Lit(param.Name), jen.Id("toString").Call(jen.Id("request").Dot(name)))

		if param.Required {
			stmts.Add(addParam)
		} else {
			stmts.If(jen.Id("request").Dot(name).Op("!=").Nil()).Block(addParam)
		}

		lastIndex++
	}

	return lastIndex
}

func (gen *goClientGenerator) generateResponse(operation *Operation, statusCode int, response spec.Response, stmts *jen.Group) {

	// Behavior of content type
	// Type
	// 	https://tools.ietf.org/search/rfc2616#section-7.2
	// Entity body
	// 	https://tools.ietf.org/search/rfc2616#section-7.2.1

	stmts.If(jen.Id("httpResponse").Dot("StatusCode").Op("==").Id(xhttp.ResolveStatusCode(statusCode))).BlockFunc(func(stmts *jen.Group) {

		hasContentType := false

		stmts.Id("contentTypeOfResponse").Op(":=").Id("extractContentType").Call(jen.Id("httpResponse").Dot("Header").Dot("Get").Call(jen.Id("contentTypeHeader")))

		if response.Schema != nil && response.Schema.Type.Contains("file") && operation.HasProduces(ContentTypesForFiles...) {

			stmts.If(jen.Id("contentTypeInList").Call(jen.Id("contentTypesForFiles"), jen.Id("contentTypeOfResponse"))).BlockFunc(func(stmts *jen.Group) {
				stmts.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID) + strconv.Itoa(statusCode) + "Response"))
				gen.generateHeaders(response.Headers, stmts)
				stmts.Id("response").Dot("Body").Op("=").Id("httpResponse").Dot("Body")
				stmts.Return(jen.Id("response"), jen.Nil())
			})

			hasContentType = true
		}

		if response.Schema != nil && (operation.HasProduce(ContentTypeApplicationJson) || operation.HasProduce(ContentTypeApplicationHalJson)) {

			stmts.If(jen.Id("contentTypeOfResponse").Op("==").Id("contentTypeApplicationJson").Op("||").Id("contentTypeOfResponse").Op("==").Id("contentTypeApplicationHalJson")).BlockFunc(func(stmts *jen.Group) {
				stmts.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID) + strconv.Itoa(statusCode) + "Response"))
				gen.generateHeaders(response.Headers, stmts)
				stmts.Id("decodeErr").Op(":=").Qual("encoding/json", "NewDecoder").Call(jen.Id("httpResponse").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("response").Dot("Body"))
				stmts.If(jen.Id("decodeErr").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Id("decodeErr")),
				)
				stmts.Return(jen.Id("response"), jen.Nil())
			}).Else().If(jen.Id("contentTypeOfResponse").Op("==").Lit("")).BlockFunc(func(stmts *jen.Group) {
				stmts.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID) + strconv.Itoa(statusCode) + "Response"))
				gen.generateHeaders(response.Headers, stmts)
				stmts.Return(jen.Id("response"), jen.Nil())
			})

			hasContentType = true
		}

		if response.Schema != nil && operation.RegexHasProduces(ContentTypeApplicationXMLPattern) {

			stmts.If(jen.Id("client").Dot("xmlMatcher").Dot("MatchString").Call(jen.Id("contentTypeOfResponse"))).BlockFunc(func(stmts *jen.Group) {
				stmts.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID) + strconv.Itoa(statusCode) + "Response"))
				gen.generateHeaders(response.Headers, stmts)
				stmts.Id("decodeErr").Op(":=").Qual("encoding/xml", "NewDecoder").Call(jen.Id("httpResponse").Dot("Body")).Dot("Decode").Call(jen.Op("&").Id("response").Dot("Body"))
				stmts.If(jen.Id("decodeErr").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Id("decodeErr")),
				)
				stmts.Return(jen.Id("response"), jen.Nil())
			}).Else().If(jen.Id("contentTypeOfResponse").Op("==").Lit("")).Block(
				jen.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID)+strconv.Itoa(statusCode)+"Response")),
				jen.Return(jen.Id("response"), jen.Nil()),
			)

			hasContentType = true
		}

		if !hasContentType {
			stmts.If(jen.Id("contentTypeOfResponse").Op("==").Lit("")).BlockFunc(func(stmts *jen.Group) {
				stmts.Id("response").Op(":=").New(jen.Id(strings.Title(operation.ID) + strconv.Itoa(statusCode) + "Response"))
				gen.generateHeaders(response.Headers, stmts)
				stmts.Return(jen.Id("response"), jen.Nil())
			})
		}

		stmts.Return(jen.Nil(), jen.Id("newNotSupportedContentType").Call(jen.Lit(415), jen.Id("contentTypeOfResponse")))

	}).Line()
}

func (gen *goClientGenerator) generateHeaders(headers map[string]spec.Header, stmts *jen.Group) {

	var headerNames []string
	for headerName := range headers {
		headerNames = append(headerNames, headerName)
	}
	sort.Strings(headerNames)

	for _, name := range headerNames {
		stmts.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("httpResponse").Dot("Header").Dot("Get").Call(jen.Lit(name)), jen.Op("&").Id("response").Dot(identifier.MakeIdentifier(strings.Title(name)))), jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		)
	}
}
