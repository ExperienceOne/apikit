package generator

import (
	"fmt"
	"github.com/ExperienceOne/apikit/generator/file"
	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/ExperienceOne/apikit/generator/stringutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var additionalServerImports = []Library{
	{Import: "github.com/go-ozzo/ozzo-routing", Alias: "routing"},
	{Import: "github.com/go-ozzo/ozzo-routing/fault"},
}

type goServerGenerator struct {
	*GoGenerator
}

func NewGoServerGenerator(spec *openapi.Spec) *goServerGenerator {

	return &goServerGenerator{
		GoGenerator: NewGoGenerator(spec),
	}
}

func (gen *goServerGenerator) Generate(path, pckg string, validators ValidatorMap, generatePrometheus bool) error {

	file := file.NewFile(pckg)
	gen.addImports(additionalServerImports, file)

	nameOfServer := strings.Title(identifier.MakeIdentifier(gen.Spec.Info().Title + "Server"))

	serverFields := []jen.Code{
		jen.Op("*").Id("Server"),
		jen.Id("Validator").Op("*").Id("Validator"),
	}

	if generatePrometheus {
		serverFields = append(serverFields, jen.Id("prometheusHandler").Op("*").Id("PrometheusHandler"))
	}

	operations := make([]*Operation, 0)
	if err := gen.WalkOperations(func(operation *Operation) error {

		parameters := make([]spec.Parameter, 0)
		parameters = append(parameters, operation.Path.Parameters...)
		parameters = append(parameters, operation.Parameters...)

		bucket := NewParameterBucket(operation.HasConsume(ContentTypeApplicationFormUrlencoded))
		bucket, err := gen.PopulateParametersBucket(bucket, parameters, operation.Security)
		if err != nil {
			return errors.Wrapf(err, "error creating parameters bucket for %s '%v'", operation.Method, operation.Path)
		}

		field, err := gen.generateOperation(operation, bucket, nameOfServer, file)
		if err != nil {
			return errors.Wrapf(err, "error generating code for %s '%s'", operation.Method, operation.Route)
		}
		serverFields = append(serverFields, field)
		operations = append(operations, operation)

		gen.generateHandler(operation, bucket, nameOfServer, generatePrometheus, file)

		return nil
	}); err != nil {
		return errors.Wrap(err, "error generating operations and handlers")
	}

	file.Type().Id(nameOfServer).Struct(serverFields...)
	gen.generateValidators(nameOfServer, validators, file)

	file.Func().Id("New" + nameOfServer).Params(jen.Id("options").Op("*").Id("ServerOpts")).Op("*").Id(nameOfServer).BlockFunc(func(stmts *jen.Group) {

		if generatePrometheus {
			stmts.Id("h").Op(":=").Id("NewPrometheusHandler").Call(jen.Nil()).Line()
			for _, op := range operations {
				route := strings.Replace(strings.Replace(op.Route, "{", "<", -1), "}", ">", -1)

				stmts.Id("h").Dot("InitMetric").Call(
					jen.Lit(route),
					jen.Lit(op.Method),
				)
			}
			stmts.Line()
		}

		wrapperArgs := []jen.Code{
			jen.Id("Server").Op(":").Id("newServer").Call(jen.Id("options")),
			jen.Id("Validator").Op(":").Id("NewValidation").Call(),
		}
		if generatePrometheus {
			wrapperArgs = append(wrapperArgs, jen.Id("prometheusHandler").Op(":").Id("h"))
		}

		stmts.Id("serverWrapper").Op(":=").Op("&").Id(nameOfServer).Values(wrapperArgs...)
		stmts.Id("serverWrapper").Dot("Server").Dot("SwaggerSpec").Op("=").Id("swagger")
		stmts.Id("serverWrapper").Dot("registerValidators").Call()
		stmts.Return(jen.Id("serverWrapper"))
	}).Line()

	file.Func().Params(jen.Id("server").Op("*").Id(nameOfServer)).Id("Start").Params(jen.Id("port").Int()).Error().BlockFunc(func(stmts *jen.Group) {

		stmts.Id("routes").Op(":=").Index().Id("RouteDescription").Values()

		if err := gen.WalkOperations(func(operation *Operation) error {
			stmts.If(jen.Id("server").Dot(stringutil.UnTitle(operation.ID + "Handler")).Op("!=").Nil()).Block(
				jen.Id("routes").Op("=").Append(jen.Id("routes"), jen.Id("server").Dot(stringutil.UnTitle(operation.ID+"Handler")).Dot("routeDescription")),
			)
			return nil
		}); err != nil {
			log.WithError(err).Error("error generating routes")
		}

		stmts.Return(jen.Id("server").Dot("Server").Dot("Start").Call(jen.Id("port"), jen.Id("routes")))
	}).Line()

	serializedSpec, err := gen.Spec.MarshalJSON()
	if err != nil {
		return err
	}
	file.Const().Id("swagger").Op("=").Lit(string(serializedSpec))

	err = file.Save(path)
	if err != nil {
		return errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}
	return nil
}

func (gen *goServerGenerator) generateValidators(nameOfServer string, validators ValidatorMap, file *file.File) {

	// group validators by regex
	regex2tags := make(map[string][]string)

	for tag, regex := range validators {
		regex2tags[regex] = append(regex2tags[regex], tag)
	}

	// sort tag names and regexs
	var regexs []string
	for regex, tags := range regex2tags {
		sort.Strings(tags)
		regex2tags[regex] = tags
		regexs = append(regexs, regex)
	}
	sort.Strings(regexs)

	file.Func().Params(jen.Id("server").Op("*").Id(nameOfServer)).Id("registerValidators").Params().BlockFunc(func(stmts *jen.Group) {
		for _, regex := range regexs {
			gen.generateValidator(regex, regex2tags[regex], stmts)
		}
	}).Line()
}

func (gen *goServerGenerator) generateValidator(regex string, tags []string, stmts *jen.Group) {

	tagIdent := identifier.MakeIdentifier(tags[0])
	callback := "callback" + strings.Title(tagIdent)

	stmts.Id(tagIdent).Op(":=").Qual("regexp", "MustCompile").Call(jen.Lit(regex))
	stmts.Id(callback).Op(":=").Func().Params(jen.Id("fl").Qual("github.com/go-playground/validator", "FieldLevel")).Id("bool").Block(
		jen.Return(jen.Id(tagIdent).Dot("MatchString").Call(jen.Id("fl").Dot("Field").Call().Dot("String").Call())),
	)
	for _, tag := range tags {
		stmts.Id("server").Dot("Validator").Dot("RegisterValidation").Call(jen.Lit(tag), jen.Id(callback))
	}
}

func (gen *goServerGenerator) generateOperation(operation *Operation, parametersBucket *ParametersBucket, nameOfServer string, file *file.File) (jen.Code, error) {

	nameOfHandler := strings.Title(operation.ID + "Handler")

	if operation.Description != "" {
		file.Comment(operation.Description)
	}
	file.Type().Id(nameOfHandler).Func().Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("request").Op("*").Id(strings.Title(operation.ID+"Request"))).Id(strings.Title(operation.ID + "Response")).Line()

	file.Type().Id(stringutil.UnTitle(nameOfHandler)+"Route").Struct(
		jen.Id("routeDescription").Id("RouteDescription"),
		jen.Id("customHandler").Id(nameOfHandler),
	).Line()

	file.Func().Params(jen.Id("server").Op("*").Id(nameOfServer)).Id("Set"+nameOfHandler).Params(jen.Id("handler").Id(nameOfHandler), jen.Id("middleware").Id("...Middleware")).Block(
		jen.Id("server").Dot(stringutil.UnTitle(nameOfHandler)).Op("=").Op("&").Id(stringutil.UnTitle(nameOfHandler)+"Route").Values(
			jen.Id("customHandler").Op(":").Id("handler"),
			jen.Id("routeDescription").Op(":").Id("RouteDescription").Values(
				jen.Id("Method").Op(":").Lit(strings.ToUpper(operation.Method)),
				jen.Id("Path").Op(":").Lit(makeMatcherForRoute(operation.Route)),
				jen.Id("Handler").Op(":").Id("server").Dot(nameOfHandler),
				jen.Id("Middleware").Op(":").Id("middleware"),
			),
		),
	).Line()

	return jen.Id(stringutil.UnTitle(nameOfHandler)).Op("*").Id(stringutil.UnTitle(nameOfHandler) + "Route"), nil
}

func (gen *goServerGenerator) generateHandler(operation *Operation, parametersBucket *ParametersBucket, nameOfServer string, generatePrometheus bool, file *file.File) {

	file.Func().Params(jen.Id("server").Op("*").Id(nameOfServer)).Id(strings.Title(operation.ID + "Handler")).Params(jen.Id("c").Op("*").Qual("github.com/go-ozzo/ozzo-routing", "Context")).Error().BlockFunc(func(stmts *jen.Group) {

		if !operation.HasValidConsumes() || !operation.HasValidProduces() {
			gen.generateNotSupported(operation, "no supported content type", false, stmts)
			return
		}

		if len(parametersBucket.Body) != 0 && operation.HasConsume(ContentTypeApplicationFormUrlencoded) {
			gen.generateNotSupported(operation, "invalid definition (both body type & application/x-www-form-urlencoded)", false, stmts)
			return
		}

		logPrefix := "wrap handler: " + operation.ID + " (" + operation.Method + ") "

		stmts.If(jen.Id("server").Dot(stringutil.UnTitle(operation.ID)+"Handler").Dot("customHandler").Op("==").Nil()).Block(
			jen.Id("server").Dot("ErrorLogger").Call(jen.Lit(logPrefix+"endpoint is not registered")),
			jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusNotFound"))),
		).Else().BlockFunc(func(stmts *jen.Group) {

			if generatePrometheus {
				stmts.Id("start").Op(":=").Qual("time", "Now").Call()
			}

			stmts.Id("request").Op(":=").New(jen.Id(operation.ID + "Request"))

			for _, param := range parametersBucket.Body {
				if param.In == "body" && operation.HasConsume(ContentTypeApplicationJson) {

					stmts.Id("contentTypeOfResponse").Op(":=").Id("extractContentType").Call(jen.Id("c").Dot("Request").Dot("Header").Dot("Get").Call(jen.Id("contentTypeHeader")))
					stmts.If(jen.Id("contentTypeOfResponse").Op("==").Id("contentTypeApplicationJson")).Block(
						jen.Id("err").Op(":=").Id("JSON").Call(jen.Id("c").Dot("Request").Dot("Body"), jen.Op("&").Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name))), jen.Lit(param.Required)),
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not decode request body of incoming request (%v)"), jen.Id("err"))),
							jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
						),
					).Else().BlockFunc(func(stmts *jen.Group) {
						logError := jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"content type of incoming request is bad (want: application/json, got: %s)"), jen.Id("contentTypeOfResponse")))
						returnUnsupported := jen.Return(jen.Id("newNotSupportedContentType").Call(jen.Lit(415), jen.Id("contentTypeOfResponse")))
						if !param.Required {
							stmts.If(jen.Id("contentTypeOfResponse").Op("!=").Lit("")).Block(logError, returnUnsupported)
						} else {
							stmts.Add(logError)
							stmts.Add(returnUnsupported)
						}
					})
				}
			}

			for _, param := range parametersBucket.Security {
				if param.Type == openapi.ApiKey.String() && param.In == openapi.Header.String() {
					stmts.Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name))).Op("=").Id("c").Dot("Request").Dot("Header").Dot("Get").Call(jen.Lit(param.Name))
				} else if param.Type == openapi.Basic.String() {
					paramName := param.Name
					if paramName == "" {
						paramName = openapi.BasicAuthHeaderName
					}
					stmts.Id("request").Dot(strings.Title(identifier.MakeIdentifier(paramName))).Op("=").Id("c").Dot("Request").Dot("Header").Dot("Get").Call(jen.Lit(paramName))
				} else {
					log.Error(fmt.Sprintf("unsupported security scheme type '%s' in '%s'", param.Type, param.In))
				}
			}

			if (parametersBucket.HasFormData && operation.HasConsume(ContentTypeMultipartFormData)) || operation.HasConsume(ContentTypeApplicationFormUrlencoded) {

				stmts.Id("contentTypeOfResponse").Op(":=").Id("extractContentType").Call(jen.Id("c").Dot("Request").Dot("Header").Dot("Get").Call(jen.Id("contentTypeHeader")))
				stmts.If(jen.Id("contentTypeOfResponse").Op("==").Id("contentTypeMultipartFormData")).BlockFunc(func(stmts *jen.Group) {

					if operation.HasConsume(ContentTypeMultipartFormData) {

						stmts.Id("formData").Op(":=").Op("&").Id("request").Dot("FormData")
						for i, param := range parametersBucket.FormDataFiles {

							file := "file" + strconv.Itoa(i)
							extractErr := "extractErr" + strconv.Itoa(i)

							stmts.List(jen.Id(file), jen.Id(extractErr)).Op(":=").Id("extractUpload").Call(jen.Lit(param.Name), jen.Id("c").Dot("Request"))
							stmts.If(jen.Id(extractErr).Op("!=").Nil()).Block(
								jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not extract upload from incoming request (error: %v)"), jen.Id(extractErr))),
								jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
							)
							if param.Required {
								file = "*" + file
							}
							stmts.Id("formData").Dot(strings.Title(identifier.MakeIdentifier(param.Name))).Op("=").Id(file)
						}

						for _, param := range parametersBucket.FormData {
							group := stmts.If(jen.Len(jen.Id("c").Dot("Request").Dot("Form").Index(jen.Lit(param.Name))).Op(">").Lit(0)).Block(
								jen.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("c").Dot("Request").Dot("Form").Index(jen.Lit(param.Name)).Index(jen.Lit(0)), jen.Op("&").Id("formData").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))), jen.Id("err").Op("!=").Nil()).Block(
									jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not convert string to specific type (error: %v)"), jen.Id("err"))),
									jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
								),
							)

							if param.Required {
								group.Else().BlockFunc(func(stmts *jen.Group) {
									stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
								})
							}
						}
					}
				}).Else().If(jen.Id("contentTypeOfResponse").Op("==").Id("contentTypeApplicationFormUrlencoded")).BlockFunc(func(stmts *jen.Group) {

					if operation.HasConsume(ContentTypeApplicationFormUrlencoded) {

						stmts.List(jen.Id("rawBody"), jen.Id("err")).Op(":=").Qual("io/ioutil", "ReadAll").Call(jen.Id("c").Dot("Request").Dot("Body"))
						stmts.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not read body of incoming request (error: %v)"), jen.Id("err"))),
							jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError"))),
						)

						stmts.List(jen.Id("queryInBody"), jen.Id("err")).Op(":=").Qual("net/url", "ParseQuery").Call(jen.String().Params(jen.Id("rawBody")))
						stmts.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not parse raw query string of incoming request (error: %v)"), jen.Id("err"))),
							jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError"))),
						)

						for _, param := range parametersBucket.FormData {
							group := stmts.If(jen.Len(jen.Id("queryInBody").Index(jen.Lit(param.Name))).Op(">").Lit(0)).Block(
								jen.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("queryInBody").Dot("Get").Call(jen.Lit(param.Name)), jen.Op("&").Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))), jen.Id("err").Op("!=").Nil()).Block(
									jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not convert string to specific type (error: %v)"), jen.Id("err"))),
									jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
								),
							)
							if param.Required {
								group.Else().BlockFunc(func(stmts *jen.Group) {
									stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
								})
							}
						}
					}

				}).Else().Block(
					jen.Return(jen.Id("newNotSupportedContentType").Call(jen.Lit(415), jen.Id("contentTypeOfResponse"))),
				)
			}

			for _, param := range parametersBucket.Path {
				// ref: https://swagger.io/docs/specification/2-0/describing-parameters/#path-parameters
				stmts.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("c").Dot("Param").Call(jen.Lit(param.Name)), jen.Op("&").Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))), jen.Id("err").Op("!=").Nil()).Block(
					jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not convert string to specific type (error: %v)"), jen.Id("err"))),
					jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
				)
			}

			for _, param := range parametersBucket.Header {
				group := stmts.If(jen.Len(jen.Id("c").Dot("Request").Dot("Header").Index(jen.Lit(param.Name))).Op(">").Lit(0)).Block(
					jen.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("c").Dot("Request").Dot("Header").Index(jen.Lit(param.Name)).Index(jen.Lit(0)), jen.Op("&").Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))), jen.Id("err").Op("!=").Nil()).Block(
						jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not convert string to specific type (error: %v)"), jen.Id("err"))),
						jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
					),
				)
				if param.Required {
					group.Else().BlockFunc(func(stmts *jen.Group) {
						stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
					})
				}
			}

			for _, param := range parametersBucket.Query {
				// ref: https://swagger.io/docs/specification/2-0/describing-parameters/#query-parameters
				parameterMemberPointer := jen.Op("&").Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))

				parameterMember := jen.Id("request").Dot(strings.Title(identifier.MakeIdentifier(param.Name)))

				group := stmts.If(jen.Len(jen.Id("c").Dot("Request").Dot("URL").Dot("Query").Call().Index(jen.Lit(param.Name))).Op(">").Lit(0)).BlockFunc(func(group *jen.Group) {
					group.If(jen.Id("err").Op(":=").Id("fromString").Call(jen.Id("c").Dot("Request").Dot("URL").Dot("Query").Call().Index(jen.Lit(param.Name)).Index(jen.Lit(0)), parameterMemberPointer), jen.Id("err").Op("!=").Nil()).Block(
						jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not convert string to specific type (error: %v)"), jen.Id("err"))),
						jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
					)

					if param.Type == "array" {
						if param.MinItems != nil {
							group.Comment("minItems validator constrain at least " + strconv.Itoa(int(*param.MinItems)) + " items")
							group.If(jen.Id("len").Call(parameterMember).Op("<").Lit(int(*param.MinItems))).Block(
								jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
							)
						}
						if param.MaxItems != nil {
							group.Comment("maxItems validator constrain at maximum " + strconv.Itoa(int(*param.MaxItems)) + " items")
							group.If(jen.Id("len").Call(parameterMember).Op(">").Lit(int(*param.MaxItems))).Block(
								jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest"))),
							)
						}
					}
				})

				if param.Required {
					group.Else().BlockFunc(func(stmts *jen.Group) {
						stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
					})
				}
			}

			stmts.List(jen.Id("validationErrors"), jen.Id("err")).Op(":=").Id("server").Dot("Validator").Dot("ValidateRequest").Call(jen.Id("request"))
			stmts.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not validate incoming request (error: %v)"), jen.Id("err"))),
				jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError"))),
			)

			stmts.If(jen.Id("validationErrors").Op("!=").Nil()).BlockFunc(func(stmts *jen.Group) {
				if response, ok := operation.Responses.StatusCodeResponses[http.StatusBadRequest]; ok {

					if response.Schema != nil && getDefinitionName(*response.Schema) == "ValidationErrors" {
						// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
						stmts.Id("c").Dot("Response").Dot("Header").Call().Index(jen.Id("contentTypeHeader")).Op("=").Index().String().Values(jen.Id("contentTypeApplicationJson"))
						stmts.Id("c").Dot("Response").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest"))
						stmts.Id("encodeErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Id("c").Dot("Response")).Dot("Encode").Call(jen.Id("validationErrors"))
						stmts.If(jen.Id("encodeErr").Op("!=").Nil()).Block(
							jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not encode validation response (error: %v)"), jen.Id("encodeErr"))),
							jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError"))),
						)
						stmts.Return(jen.Nil())
					} else {
						stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
					}

				} else {
					stmts.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusBadRequest")))
				}
			})

			stmts.Id("response").Op(":=").Id("server").Dot(stringutil.UnTitle(operation.ID)+"Handler").Dot("customHandler").Call(jen.Id("c").Dot("Request").Dot("Context").Call(), jen.Id("request"))
			stmts.If(jen.Id("response").Op("==").Nil()).Block(
				jen.Id("server").Dot("ErrorLogger").Call(jen.Lit(logPrefix+"received a nil response object")),
				jen.Return(jen.Id("NewHTTPStatusCodeError").Call(jen.Qual("net/http", "StatusInternalServerError"))),
			)

			stmts.If(jen.Id("err").Op(":=").Id("response").Dot("write").Call(jen.Id("c").Dot("Response")), jen.Id("err").Op("!=").Nil()).Block(
				jen.Id("server").Dot("ErrorLogger").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(logPrefix+"could not send response (error: %v)"), jen.Id("err"))),
				jen.Return(jen.Id("err")),
			)

			if generatePrometheus {
				stmts.Id("routeDescription").Op(":=").Id("server").Dot(stringutil.UnTitle(operation.ID) + "Handler").Dot("routeDescription")
				stmts.Id("server").Dot("prometheusHandler").Dot("HandleRequest").Call(
					jen.Id("routeDescription").Dot("Path"),
					jen.Id("routeDescription").Dot("Method"),
					jen.Id("response").Dot("StatusCode").Call(),
					jen.Qual("time", "Since").Call(jen.Id("start")),
				)
			}
		})

		stmts.Return(jen.Nil())
	}).Line()
}
