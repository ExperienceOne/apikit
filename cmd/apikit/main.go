package main

import (
	"os"

	"github.com/ExperienceOne/apikit/generator"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/ExperienceOne/apikit/internal/framework/version"

	openapierror "github.com/go-openapi/errors"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	cmdGenerate string = "generate"
	cmdProject  string = "project"
	cmdValidate string = "validate"
	cmdHandler  string = "handlers"
	cmdService  string = "service"
	cmdVersion  string = "version"

	flagDebug              string = "debug"
	flagDebugAlias         string = "d"
	flagGenerateOnlyClient string = "only-client"
	flagGenerateOnlyServer string = "only-server"
	flagGenerateMock       string = "mocked"
	flagGeneratePrometheus string = "prometheus"
)

func main() {

	app := cli.NewApp()
	app.Name = "apikit"
	app.Description = "apikit generates server and client Go code based on OpenAPIv2 (Swagger) definitions"
	app.Usage = "apikit <project|generate|validate|handlers|service|version>"
	app.Version = version.GitTag

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagDebug + ", " + flagDebugAlias,
			Usage: "activate debug mode",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        cmdProject,
			Description: "creates the common Go project structure",
			Usage:       "apikit project <dest.dir> <path/of/package>",
			Action: func(c *cli.Context) error {
				return generator.GoProject(c.Args().Get(0), c.Args().Get(1))
			},
		},
		{
			Name:        cmdGenerate,
			Description: "creates or updates generated code based on an OpenAPIv2 (Swagger) definition",
			Usage:       "apikit generate <api.yaml> <dest> <package>",
			Action: func(ctx *cli.Context) error {

				generatePrometheus := ctx.Bool(flagGeneratePrometheus)

				if ctx.Bool(flagGenerateOnlyClient) {
					if ctx.Bool(flagGenerateMock) {
						return GenerateAction(generator.NewGoClientAPIMockGenerator, generatePrometheus, ctx)
					} else {
						return GenerateAction(generator.NewGoClientAPIGenerator, generatePrometheus, ctx)
					}
				}
				if ctx.Bool(flagGenerateOnlyServer) {
					if ctx.Bool(flagGenerateMock) {
						return GenerateAction(generator.NewGoServerAPIMockGenerator, generatePrometheus, ctx)
					} else {
						return GenerateAction(generator.NewGoServerAPIGenerator, generatePrometheus, ctx)
					}
				}
				if ctx.Bool(flagGenerateMock) {
					return GenerateAction(generator.NewGoAPIMockGenerator, generatePrometheus, ctx)
				} else {
					return GenerateAction(generator.NewGoAPIGenerator, generatePrometheus, ctx)
				}
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  flagGenerateOnlyClient,
					Usage: "generate client code only",
				},
				cli.BoolFlag{
					Name:  flagGenerateOnlyServer,
					Usage: "generate server code only",
				},
				cli.BoolFlag{
					Name:  flagGenerateMock,
					Usage: "generate mock",
				},
				cli.BoolFlag{
					Name:  flagGeneratePrometheus,
					Usage: "generate Prometheus handlers",
				},
			},
		},
		{
			Name:        cmdValidate,
			Description: "validates an OpenAPIv2 (Swagger) definition",
			Usage:       "apikit validate <api.yaml>",
			Action:      ValidateAction,
		},
		{
			Name:        cmdHandler,
			Description: "creates stubs for the API endpoint handlers of the OpenAPIv2 (Swagger) definition",
			Usage:       "apikit handlers <api.yaml> <dest.go> <package> <api/package/path>",
			Action: func(ctx *cli.Context) error {
				return GenerateHandlersAction(ctx)
			},
		},
		{
			Name:        cmdService,
			Description: "creates service stub for the tagged API endpoint handlers of the OpenAPIv2 (Swagger) definition",
			Usage:       "apikit service <api.yaml> <dest.go> <package> <tag> <api/package/path>",
			Action: func(ctx *cli.Context) error {
				return GenerateServiceAction(ctx)
			},
		},
		{
			Name:        cmdVersion,
			Description: "prints version information",
			Usage:       "apikit version",
			Action: func(ctx *cli.Context) error {
				return version.ApikitVersion().PrintTable()
			},
		},
	}

	args := os.Args
	if len(os.Args) > 1 && (args[1] == "--debug" || args[1] == "-d") {
		args = args[2:]
	} else if len(args) > 0 {
		args = args[1:]
	}

	if len(args) == 0 {
		cli.ShowAppHelpAndExit(cli.NewContext(app, nil, nil), 1)
	}

	command := args[0]
	if command != cmdProject && command != cmdGenerate && command != cmdValidate && command != cmdHandler && command != cmdService && command != cmdVersion {
		cli.ShowAppHelpAndExit(cli.NewContext(app, nil, nil), 1)
	}

	if command == cmdProject && len(args) != 3 {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdProject, 1)
	}

	if command == cmdGenerate && (len(args) < 4 || len(args) > 5) {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdGenerate, 1)
	}

	if command == cmdValidate && len(args) != 2 {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdValidate, 1)
	}

	if command == cmdHandler && len(args) != 5 {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdHandler, 1)
	}

	if command == cmdService && len(args) != 6 {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdService, 1)
	}

	if command == cmdVersion && len(args) != 1 {
		cli.ShowCommandHelpAndExit(cli.NewContext(app, nil, nil), cmdVersion, 1)
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("error running apikit")
	}
}

func GenerateAction(constructor func(spec *openapi.Spec) generator.Generator, generatePrometheus bool, ctx *cli.Context) error {

	if ctx.GlobalBool(flagDebug) {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode activated")
	}

	specFile, dest, pkg := ctx.Args().Get(0), ctx.Args().Get(1), ctx.Args().Get(2)

	spec, err := openapi.NewOpenApiSpecFromFile(specFile)
	if err != nil {
		return errors.Wrapf(err, "failed to load swagger file '%s'", specFile)
	}

	if err := constructor(spec).Generate(dest, pkg, generatePrometheus); err != nil {
		return errors.Wrap(err, "failed to generate code")
	}

	return nil
}

func GenerateHandlersAction(ctx *cli.Context) error {

	if ctx.GlobalBool(flagDebug) {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode activated")
	}

	specFile, dest, pkg, serverPkg := ctx.Args().Get(0), ctx.Args().Get(1), ctx.Args().Get(2), ctx.Args().Get(3)

	spec, err := openapi.NewOpenApiSpecFromFile(specFile)
	if err != nil {
		return errors.Wrapf(err, "failed to load swagger file '%s'", specFile)
	}

	if err := generator.NewGoHandlersGenerator(spec).Generate(dest, pkg, serverPkg); err != nil {
		return errors.Wrap(err, "failed to generate handlers")
	}

	return nil
}

func GenerateServiceAction(ctx *cli.Context) error {

	if ctx.GlobalBool(flagDebug) {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode activated")
	}

	specFile, dest, pkg, tag, serverPkg := ctx.Args().Get(0), ctx.Args().Get(1), ctx.Args().Get(2), ctx.Args().Get(3), ctx.Args().Get(4)

	spec, err := openapi.NewOpenApiSpecFromFile(specFile)
	if err != nil {
		return errors.Wrapf(err, "failed to load swagger file '%s'", specFile)
	}

	if err := generator.NewGoServiceGenerator(spec).Generate(dest, pkg, tag, serverPkg); err != nil {
		return errors.Wrap(err, "failed to generate service")
	}

	return nil
}

func ValidateAction(ctx *cli.Context) error {

	if ctx.GlobalBool(flagDebug) {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode activated")
	}

	specFile := ctx.Args().Get(0)
	spec, err := openapi.NewOpenApiSpecFromFile(specFile)
	if err != nil {
		return errors.Wrapf(err, "failed to load swagger file '%s'", specFile)
	}

	if err := spec.Validate(); err != nil {
		if err, ok := err.(*openapierror.CompositeError); ok {
			for _, subErr := range err.Errors {
				log.Error(subErr)
			}
		} else {
			log.Error(err)
		}
	} else {
		log.Info("Swagger definition is valid")
	}

	return nil
}
