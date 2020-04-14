package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ExperienceOne/apikit/generator/openapi"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type goServicesGenerator struct {
	*GoGenerator
}

func NewGoServiceGenerator(spec *openapi.Spec) *goServicesGenerator {

	return &goServicesGenerator{
		GoGenerator: NewGoGenerator(spec),
	}
}

func (gen *goServicesGenerator) Generate(path, pckg, tag, serverPkg string) error {

	log.Info(fmt.Sprintf("start generating service for '%s'", tag))
	defer func() {
		log.Info(fmt.Sprintf("finished generating service for '%s'", tag))
	}()

	filteredOperations := make([]*Operation, 0)
	if err := gen.WalkOperations(func(operation *Operation) error {
		if len(operation.Tags) == 0 {
			return nil
		}
		t := operation.Tags[0]
		if t == tag {
			filteredOperations = append(filteredOperations, operation)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "error generating service")
	}

	if len(filteredOperations) == 0 {
		return errors.Errorf("no operations for tag '%s'", tag)
	}

	file := jen.NewFile(pckg)
	file.Type().Id(strings.Title(tag) + "Service").Struct()
	for _, operation := range filteredOperations {
		log.Info(fmt.Sprintf("generate service handler for %s '%s'", operation.Method, operation.ID))
		gen.generateHandler(tag, operation, serverPkg, file)
	}
	err := file.Save(path)
	if err != nil {
		return errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}
	return nil
}

func (gen *goServicesGenerator) generateHandler(service string, operation *Operation, serverPkg string, file *jen.File) {

	file.Func().Params(jen.Id("service").Op("*").Id(strings.Title(service)+"Service")).Id(strings.Title(operation.ID)).Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("request").Op("*").Qual(serverPkg, strings.Title(operation.ID)+"Request")).Qual(serverPkg, strings.Title(operation.ID)+"Response").BlockFunc(func(stmts *jen.Group) {
		response, status, _ := operation.SuccessResponse()
		if response == nil {
			stmts.Return(jen.Nil())
		} else {
			stmts.Return(jen.Op("&").Qual(serverPkg, strings.Title(operation.ID)+strconv.Itoa(status)+"Response").Values())
		}
	}).Line()
}
