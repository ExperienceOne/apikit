package generator

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/ExperienceOne/apikit/generator/file"
	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GoGenerator struct {
	Spec *openapi.Spec
}

var methods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
}

const (
	ContentTypeApplicationJson           string = "application/json"
	ContentTypeApplicationHalJson        string = "application/hal+json"
	ContentTypeMultipartFormData         string = "multipart/form-data"
	ContentTypeImagePng                  string = "image/png"
	ContentTypeImageJpeg                 string = "image/jpeg"
	ContentTypeImageTiff                 string = "image/tiff"
	ContentTypeImageWebp                 string = "image/webp"
	ContentTypeImageGif                  string = "image/gif"
	ContentTypeImageSvgXml               string = "image/svg+xml"
	ContentTypeImageXIcon                string = "image/x-icon"
	ContentTypeTextPlain                 string = "text/plain; charset=utf-8"
	ContentTypeTextHTML                  string = "text/html"
	ContentTypeApplicationFormUrlencoded string = "application/x-www-form-urlencoded"
	ContentTypeApplicationPDF            string = "application/pdf"
	ContentTypeApplicationXMLPattern     string = `^application\/(.+)xml$`
	ContentTypeApplicationOctetStream string = "application/octet-stream"
)

var ContentTypesForFiles []string = []string{
	ContentTypeApplicationJson,
	ContentTypeImagePng,
	ContentTypeImageJpeg,
	ContentTypeImageTiff,
	ContentTypeImageWebp,
	ContentTypeImageSvgXml,
	ContentTypeImageGif,
	ContentTypeImageTiff,
	ContentTypeImageXIcon,
	ContentTypeApplicationPDF,
	ContentTypeApplicationOctetStream,
}

var supportedProduces []string = []string{
	ContentTypeApplicationJson,
	ContentTypeApplicationHalJson,
	ContentTypeApplicationFormUrlencoded,
	ContentTypeMultipartFormData,
	ContentTypeTextPlain,
	// content types for images
	ContentTypeImagePng,
	ContentTypeImageJpeg,
	ContentTypeImageTiff,
	ContentTypeImageWebp,
	ContentTypeImageSvgXml,
	ContentTypeImageGif,
	ContentTypeImageTiff,
	ContentTypeImageXIcon,
	// content types for additional files
	ContentTypeApplicationPDF,
	ContentTypeApplicationXMLPattern,
	ContentTypeApplicationOctetStream,
}

var supportedConsumes []string = []string{
	ContentTypeApplicationJson,
	ContentTypeApplicationHalJson,
	ContentTypeApplicationFormUrlencoded,
	ContentTypeMultipartFormData,
	ContentTypeTextPlain,
	// content types for images
	ContentTypeImagePng,
	ContentTypeImageJpeg,
	ContentTypeImageTiff,
	ContentTypeImageWebp,
	ContentTypeImageSvgXml,
	ContentTypeImageGif,
	ContentTypeImageTiff,
	ContentTypeImageXIcon,
	// content types for additional files
	ContentTypeApplicationPDF,
	ContentTypeApplicationXMLPattern,
}

type Generator interface {
	Generate(path string, pkg string, generatePrometheus bool) error
}

func NewGoGenerator(spec *openapi.Spec) *GoGenerator {

	return &GoGenerator{Spec: spec}
}

type Library struct {
	Import string
	Alias  string
}

func (gen *GoGenerator) addImports(imports []Library, file *file.File) {

	for _, imprt := range imports {
		if imprt.Alias != "" {
			file.ImportAlias(imprt.Import, imprt.Alias)
		} else {
			file.ImportName(imprt.Import, "")
		}
	}
}

func (gen *GoGenerator) WalkOperations(handler func(operation *Operation) error) error {

	paths := gen.Spec.Paths()

	i := 0
	routes := make([]string, len(paths))
	for route := range paths {
		routes[i] = route
		i++
	}
	sort.Strings(routes)

	for _, route := range routes {
		path := paths[route]
		rPath := reflect.ValueOf(path)
		for _, method := range methods {
			methodName := strings.Title(strings.ToLower(method))
			if !rPath.FieldByName(methodName).IsNil() {

				operation := rPath.FieldByName(methodName).Interface().(*spec.Operation)

				operationID, err := identifier.ValidateAndCleanOperationsID(operation.ID)
				if err != nil {
					return errors.Wrapf(err, "error validating operation ID of path %s '%s'", method, route)
				}
				operation.ID = operationID

				o := &Operation{
					Method:    method,
					Route:     route,
					Path:      &path,
					Operation: operation,
				}

				if len(operation.Consumes) != 0 {
					o.Consumes = filterContentTypes(operation.Consumes, supportedConsumes)
				} else if len(gen.Spec.GlobalConsumes()) != 0 {
					o.Consumes = filterContentTypes(gen.Spec.GlobalConsumes(), supportedConsumes)
				} else {
					o.Consumes = []string{ContentTypeApplicationJson}
				}

				if len(operation.Produces) != 0 {
					o.Produces = filterContentTypes(operation.Produces, supportedProduces)
				} else if len(gen.Spec.GlobalProduces()) != 0 {
					o.Produces = filterContentTypes(gen.Spec.GlobalProduces(), supportedProduces)
				} else {
					o.Produces = []string{ContentTypeApplicationJson}
				}

				if err := handler(o); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (gen *GoGenerator) ConvertMapToSecuritySchemeSlice(security []map[string][]string) []spec.SecurityScheme {

	var securityParams []spec.SecurityScheme
	for _, security := range security {
		for securityScheme := range security {
			securityParams = append(securityParams, *gen.Spec.SecurityScheme(securityScheme))
		}
	}
	return securityParams
}

type ParametersBucket struct {
	Path          []*spec.Parameter
	Query         []*spec.Parameter
	Header        []*spec.Parameter
	FormData      []*spec.Parameter
	FormDataFiles []*spec.Parameter
	Security      []spec.SecurityScheme
	Body          []*spec.Parameter
	HasFormData   bool
	HasBody       bool
	HasURLEncoded bool
}

func NewParameterBucket(HasURLEncoded bool) *ParametersBucket {
	return &ParametersBucket{
		Path:          make([]*spec.Parameter, 0),
		Query:         make([]*spec.Parameter, 0),
		Header:        make([]*spec.Parameter, 0),
		FormData:      make([]*spec.Parameter, 0),
		FormDataFiles: make([]*spec.Parameter, 0),
		Body:          make([]*spec.Parameter, 0),
		HasURLEncoded: HasURLEncoded,
	}
}

func (gen *GoGenerator) resolveParam(param *spec.Parameter) (*spec.Parameter, error) {

	if !param.Ref.GetPointer().IsEmpty() {
		var err error
		param, err = spec.ResolveParameter(gen.Spec.Spec, param.Ref)
		if err != nil {
			return nil, err
		} else {
			return param, nil
		}
	}

	return param, nil
}

func (gen *GoGenerator) PopulateParametersBucket(bucket *ParametersBucket, parameters []spec.Parameter, security []map[string][]string) (*ParametersBucket, error) {

	// operation level security schemes replace the global schemes (no merging, see spec!)
	var securitySchemes []spec.SecurityScheme
	if len(security) > 0 {
		securitySchemes = gen.ConvertMapToSecuritySchemeSlice(security)
	} else {
		securitySchemes = gen.ConvertMapToSecuritySchemeSlice(gen.Spec.GlobalSecurities())
	}

	// remove double header fields
	for _, security := range securitySchemes {
		if security.In == openapi.Header.String() {
			alreadyRegistered := false
			for _, s := range bucket.Security {
				if s.In == openapi.Header.String() && s.Name == security.Name {
					alreadyRegistered = true
					break
				}
			}
			if !alreadyRegistered {
				bucket.Security = append(bucket.Security, security)
			}
		} else {
			bucket.Security = append(bucket.Security, security)
		}
	}

	sortParamFunc := func(param spec.Parameter) error {

		resolvedParam, err := gen.resolveParam(&param)
		if err != nil {
			return errors.Wrapf(err, "error resolving param ref '%s'", param.Ref.GetURL())
		}

		if resolvedParam.In == openapi.Query.String() {
			bucket.Query = append(bucket.Query, resolvedParam)
		} else if resolvedParam.In == openapi.Path.String() {
			bucket.Path = append(bucket.Path, resolvedParam)
		} else if resolvedParam.In == openapi.Header.String() {
			bucket.Header = append(bucket.Header, resolvedParam)
		} else if resolvedParam.In == openapi.Body.String() {
			bucket.HasBody = true
			bucket.Body = append(bucket.Body, resolvedParam)
		} else if resolvedParam.In == openapi.FormData.String() && bucket.HasURLEncoded {
			bucket.FormData = append(bucket.FormData, resolvedParam)
		} else if resolvedParam.In == openapi.FormData.String() && bucket.HasBody {
			bucket.Body = append(bucket.Body, resolvedParam)
		} else if resolvedParam.In == openapi.FormData.String() && resolvedParam.Type == "file" {
			bucket.HasFormData = true
			bucket.FormDataFiles = append(bucket.FormDataFiles, resolvedParam)
		} else if resolvedParam.In == openapi.FormData.String() {
			bucket.HasFormData = true
			bucket.FormData = append(bucket.FormData, resolvedParam)
		} else {
			return errors.Errorf("unknown 'in' property '%s' of param '%s'", resolvedParam.In, resolvedParam.Name)
		}

		return nil
	}

	for _, param := range parameters {
		err := sortParamFunc(param)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}

func (gen *GoGenerator) generateNotSupported(operation *Operation, message string, returnWithNil bool, stmts *jen.Group) {

	log.Info(fmt.Sprintf("generating unsupported content type handler for operation '%s': %s", operation.ID, message))

	notSupported := jen.Id("newNotSupportedContentType").Call(jen.Lit(415), jen.Lit(message))

	if returnWithNil {
		stmts.Return(jen.List(jen.Nil(), notSupported))
	} else {
		stmts.Return(notSupported)
	}
}
