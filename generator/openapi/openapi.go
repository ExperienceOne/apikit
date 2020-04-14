package openapi

import (
	"path/filepath"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

type Spec struct {
	Spec *spec.Swagger
	doc  *loads.Document
}

func NewOpenApiSpecFromFile(path string) (*Spec, error) {
	doc, err := fileToDocument(path)
	if err != nil {
		return nil, err
	}
	return &Spec{
		Spec: doc.Spec(),
		doc:  doc,
	}, nil
}

func (oas *Spec) Validate() error {
	if err := validate.Spec(oas.doc, strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (oas *Spec) Info() *spec.InfoProps {
	return &oas.Spec.Info.InfoProps
}

// A list of MIME types the APIs on this resource can consume.
// This is global to all APIs but can be overridden on specific API calls.
func (oas *Spec) GlobalConsumes() []string {
	return oas.Spec.Consumes
}

//	A list of MIME types the APIs on this resource can produce.
//  This is global to all APIs but can be overridden on specific API calls.
func (oas *Spec) GlobalProduces() []string {
	return oas.Spec.Produces
}

func (oas *Spec) Definitions() spec.Definitions {
	return oas.Spec.Definitions
}

func (oas *Spec) Parameters() map[string]spec.Parameter {
	return oas.Spec.Parameters
}

func (oas *Spec) Paths() map[string]spec.PathItem {

	return oas.Spec.Paths.Paths
}

func (oas *Spec) SecurityScheme(name string) *spec.SecurityScheme {

	return oas.Spec.SecurityDefinitions[name]
}

func (oas *Spec) GlobalSecurities() []map[string][]string {

	return oas.Spec.Security
}

func (oas *Spec) MarshalJSON() ([]byte, error) {

	return oas.Spec.MarshalJSON()
}

const (
	jsonExt string = ".json"
	yamlExt string = ".yaml"
)

func fileToDocument(f string) (*loads.Document, error) {

	var doc *loads.Document
	var err error

	switch filepath.Ext(f) {
	case jsonExt:
		doc, err = loads.JSONSpec(f)
	case yamlExt:
		doc, err = loads.Spec(f)
	default:
		return nil, errors.Errorf("error loading swagger from '%s', extension isn't supported by generator", f)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error loading swagger from '%s'", f)
	}

	analysisSpec := analysis.New(doc.Spec())
	opts := analysis.FlattenOpts{
		Spec:     analysisSpec,
		BasePath: f,
		Minimal:  true,
	}

	if err := analysis.Flatten(opts); err != nil {
		return nil, errors.Wrap(err, "failed to merge swagger specifications")
	}

	return doc, nil
}
