package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ExperienceOne/apikit/generator/templates"

	"github.com/pkg/errors"
)

const (
	permission = 0775
)

func GoProject(dir, pckg string) error {

	err := os.MkdirAll(dir, permission)
	if err != nil {
		return errors.Wrap(err, "error creating project directory")
	}

	name := filepath.Base(pckg)

	err = os.MkdirAll(filepath.Join(dir, "cmd", name), permission)
	if err != nil {
		return errors.Wrap(err, "error creating command directory")
	}

	err = writeToFile(filepath.Join(dir, "cmd", name, "main.go"), fmt.Sprintf(templates.MainTpl, pckg, name, name))
	if err != nil {
		return errors.Wrap(err, "error writing main.go")
	}

	err = os.MkdirAll(filepath.Join(dir, "doc"), permission)
	if err != nil {
		return errors.Wrap(err, "error creating doc directory")
	}

	err = os.MkdirAll(filepath.Join(dir, "cfg"), permission)
	if err != nil {
		return errors.Wrap(err, "error creating config directory")
	}

	err = os.MkdirAll(filepath.Join(dir, "api"), permission)
	if err != nil {
		return errors.Wrap(err, "error creating api directory")
	}

	err = writeToFile(filepath.Join(dir, ".gitignore"), templates.GitignoreTpl)
	if err != nil {
		return errors.Wrap(err, "error writing .gitignore")
	}

	return nil
}

func writeToFile(name, content string) error {

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
