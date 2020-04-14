package generator

import (
	"strconv"
	"strings"

	"github.com/ExperienceOne/apikit/generator/identifier"
	"github.com/ExperienceOne/apikit/generator/openapi"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type goHandlersGenerator struct {
	*GoGenerator
}

func NewGoHandlersGenerator(spec *openapi.Spec) *goHandlersGenerator {

	return &goHandlersGenerator{
		GoGenerator: NewGoGenerator(spec),
	}
}

func (gen *goHandlersGenerator) Generate(path, pckg, serverPkg string) error {

	file := jen.NewFile(pckg)

	file.Func().Id("RegisterHandlers").Params(jen.Id("server").Op("*").Qual(serverPkg, strings.Title(identifier.MakeIdentifier(gen.Spec.Info().Title+"Server")))).BlockFunc(func(stmts *jen.Group) {
		if err := gen.WalkOperations(func(operation *Operation) error {
			gen.generateHandler(operation, serverPkg, stmts, file)
			return nil
		}); err != nil {
			log.WithError(err).Error("error generating handlers")
		}
	}).Line()

	err := file.Save(path)
	if err != nil {
		return errors.Wrapf(err, "error writing generated code to file '%s'", path)
	}
	return nil
}

func (gen *goHandlersGenerator) generateHandler(operation *Operation, serverPkg string, stmts *jen.Group, file *jen.File) {

	file.Func().Id(strings.Title(operation.ID)).Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("request").Op("*").Qual(serverPkg, strings.Title(operation.ID)+"Request")).Qual(serverPkg, strings.Title(operation.ID)+"Response").BlockFunc(func(stmts *jen.Group) {
		response, status, _ := operation.SuccessResponse()
		if response == nil {
			stmts.Return(jen.Nil())
		} else {
			stmts.Return(jen.Op("&").Qual(serverPkg, strings.Title(operation.ID)+strconv.Itoa(status)+"Response").Values())
		}
	}).Line()

	stmts.Id("server").Dot("Set" + strings.Title(operation.ID) + "Handler").Call(jen.Id(strings.Title(operation.ID)))
}
