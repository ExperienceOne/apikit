package generator

import (
	"github.com/spf13/viper"
	"github.com/vektra/mockery/v2/cmd"
	"github.com/vektra/mockery/v2/pkg/config"
	"strings"
)

type mockGoClientGenerator struct{}

func NewMockGoClientGenerator() *mockGoClientGenerator {
	return &mockGoClientGenerator{}
}

func (gen *mockGoClientGenerator) Generate(clientName, path string) error {
	mockery, err := cmd.GetRootAppFromViper(viper.GetViper())
	if err != nil {
		return err
	}
	mockery.Config.InPackage = true
	mockery.Config.Name = strings.Title(clientName)
	mockery.Config.Dir = path
	config.SemVer = "v2.9.0"

	return mockery.Run()
}
