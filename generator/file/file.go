package file

import (
	"github.com/dave/jennifer/jen"
)

type File struct {
	*jen.File
	types map[string]bool
}

func NewFile(packageName string) *File {

	return &File{
		File:  jen.NewFile(packageName),
		types: make(map[string]bool),
	}
}

func (file *File) HasType(name string) bool {

	return file.types[name]
}

func (file *File) AddType(name string) *jen.Statement {

	file.types[name] = true
	return file.Type().Id(name)
}
