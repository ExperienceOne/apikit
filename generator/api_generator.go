package generator

import (
	"github.com/spf13/viper"
	"github.com/vektra/mockery/v2/cmd"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ExperienceOne/apikit/framework"
	"github.com/ExperienceOne/apikit/generator/openapi"
	"github.com/pkg/errors"
)

const (
	frameworkFile = "framework.go"
	typesFile     = "types.go"
	serverFile    = "server.go"
	clientFile    = "client.go"
)

var (
	frameworkIdentifiers = []string{
		"SetRequestHeadersFromContext",
		"NewHttpContextWrapper",
		"ExtractContentType",
		"ContentTypeInList",
		"ContentTypeHeader",
		"ContentTypeApplicationJson",
		"ContentTypeApplicationHalJson",
		"ContentTypeMultipartFormData",
		"NewNotSupportedContentType",
		"ContentTypeApplicationFormUrlencoded",
		"ToString",
		"NewErrOnUnknownResponseCode",
		"NewErrUnknownResponse",
		"ExtractUpload",
		"FromString",
		"NewServer",
		"NewHttpClientWrapper",
		"HttpClientWrapper",
		"NewRequestObjectIsNilError",
		"ServeJson",
		"ServeHalJson",
		"HttpCodeError",
		"XHTTPError",
		"NewJsonHTTPError",
		"JsonHTTPError",
	}

	frameworkPackages = []string{
		"hooks",
		"middleware",
		"parameter",
		"roundtripper",
		"validation",
		"version",
		"xclient",
		"xhttp",
		"xhttperror",
		"xserver",
	}
)

type ApiGenerator struct {
	GenerateMocks bool
	goTypesGenerator      *goTypesGenerator
	goServerGenerator     *goServerGenerator
	goClientGenerator     *goClientGenerator
}

func NewGoAPIGenerator(spec *openapi.Spec) Generator {

	return &ApiGenerator{
		GenerateMocks:         true,
		goTypesGenerator:      NewGoTypesGenerator(spec),
		goServerGenerator:     NewGoServerGenerator(spec),
		goClientGenerator:     NewGoClientGenerator(spec),
	}
}

func (gen *ApiGenerator) Generate(path, pkg string, generatePrometheus bool) error {

	return gen.generate(path, pkg, true, true, false, generatePrometheus)
}

func (gen *ApiGenerator) generate(path, pkg string, client, server, mock, generatePrometheus bool) error {

	fwCode := framework.Code
	if client && !server {
		fwCode = framework.ClientCode
	} else if server && !client {
		fwCode = framework.ServerCode
	}

	fw, err := framework.FromBytes(fwCode)
	if err != nil {
		return errors.Wrap(err, "error loading framework code")
	}

	fw.RenamePackage(pkg)

	err = fw.MakePrivate(frameworkIdentifiers)
	if err != nil {
		return errors.Wrap(err, "error making frameworks identifiers private")
	}

	err = fw.RenameTypes(frameworkPackages, "")
	if err != nil {
		return errors.Wrap(err, "error collapsing framework packages")
	}

	source, err := fw.Bytes()
	if err != nil {
		return errors.Wrap(err, "error getting framework source code")
	}

	if err := ioutil.WriteFile(filepath.Join(path, frameworkFile), source, 0644); err != nil {
		return errors.Wrap(err, "error persisting framework code")
	}

	validators, err := gen.goTypesGenerator.Generate(filepath.Join(path, typesFile), pkg)
	if err != nil {
		return errors.Wrap(err, "error generating types")
	}

	if server {
		if mock {
			return errors.New("server mock is not implemented yet")
		} else {
			err = gen.goServerGenerator.Generate(filepath.Join(path, serverFile), pkg, validators, generatePrometheus)
			if err != nil {
				return errors.Wrap(err, "error generating server")
			}
		}
	}

	if client {
		err = gen.goClientGenerator.Generate(filepath.Join(path, clientFile), pkg, generatePrometheus)
		if err != nil {
			return errors.Wrap(err, "error generating client")
		}

		if gen.GenerateMocks {
			mockery, err := cmd.GetRootAppFromViper(viper.GetViper())
			if err != nil {
				return err
			}
			mockery.Config.InPackage = true
			mockery.Config.Name = strings.Title(gen.goClientGenerator.ClientName())
			mockery.Config.Dir = path

			err = mockery.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type ServerApiGenerator struct {
	ApiGenerator
}

func NewGoServerAPIGenerator(spec *openapi.Spec) Generator {

	return &ServerApiGenerator{
		ApiGenerator: ApiGenerator{
			goServerGenerator: NewGoServerGenerator(spec),
			goTypesGenerator:  NewGoTypesGenerator(spec),
		},
	}
}

func (gen *ServerApiGenerator) Generate(path, pkg string, generatePrometheus bool) error {

	return gen.generate(path, pkg, false, true, false, generatePrometheus)
}

type ClientApiGenerator struct {
	ApiGenerator
}

func NewGoClientAPIGenerator(spec *openapi.Spec) Generator {

	return &ClientApiGenerator{
		ApiGenerator: ApiGenerator{
			goClientGenerator: NewGoClientGenerator(spec),
			goTypesGenerator:  NewGoTypesGenerator(spec),
		},
	}
}

func (gen *ClientApiGenerator) Generate(path, pkg string, generatePrometheus bool) error {
	return gen.generate(path, pkg, true, false, false, generatePrometheus)
}
