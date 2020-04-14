package generator

import (
	"github.com/ExperienceOne/apikit/generator/stringutil"

	"github.com/go-openapi/spec"
)

type Operation struct {
	*spec.Operation
	Method   string
	Route    string
	Path     *spec.PathItem
	Produces []string
	Consumes []string
}

func (o *Operation) HasConsume(c string) bool {
	return stringutil.InStringSlice(o.Consumes, c)
}

func (o *Operation) RegexHasProduces(pattern string) bool {
	return regexContains(o.Produces, pattern)
}

func (o *Operation) HasConsumes(cs ...string) bool {
	for _, c := range cs {
		if o.HasConsume(c) {
			return true
		}
	}
	return false
}

func (o *Operation) HasProduce(p string) bool {
	return stringutil.InStringSlice(o.Produces, p)
}

func (o *Operation) HasProduces(ps ...string) bool {
	for _, p := range ps {
		if o.HasProduce(p) {
			return true
		}
	}
	return false
}

func (o *Operation) HasValidConsumes() bool {
	return len(o.Consumes) != 0
}

func (o *Operation) HasValidProduces() bool {
	return len(o.Produces) != 0
}
