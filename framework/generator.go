package framework

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ExperienceOne/apikit/generator/stringutil"

	"github.com/fatih/astrewrite"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
)

var (
	packageDeclaration = ([]byte)("\n\tpackage framework\n")
)

type Generator struct {
	code  *ast.File
	files *token.FileSet
}

// Creates new framework source from directory
func FromDirectory(dir string, excludedPackages []string) (*Generator, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if len(excludedPackages) > 0 && containsPath(excludedPackages, path) {
			return nil
		}

		if info.IsDir() || strings.Contains(path, "_test.go") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to read framework files from source directory '%s'", dir)
	}

	fileSet := token.NewFileSet()
	dst, err := parser.ParseFile(fileSet, "", packageDeclaration, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse package declaration")
	}

	for _, file := range files {
		src, err := parser.ParseFile(fileSet, file, nil, 0)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse file '%s'", file)
		}

		missingImports := make([]*ast.ImportSpec, 0)
		for _, importSpec := range src.Imports {
			if strings.Contains(importSpec.Path.Value, "apikit") {
				continue
			}

			contains := func(arr []*ast.ImportSpec, imp *ast.ImportSpec) bool {
				for _, item := range arr {
					if item.Path.Value == imp.Path.Value {
						return true
					}
				}
				return false
			}

			if !contains(dst.Imports, importSpec) {
				missingImports = append(missingImports, importSpec)
			}
		}

		for _, importSpecs := range astutil.Imports(fileSet, src) {
			for _, importSpec := range importSpecs {
				if ok := astutil.DeleteImport(fileSet, src, strings.Replace(importSpec.Path.Value, "\"", "", -1)); !ok {
					return nil, errors.Errorf("failed to delete import '%s' of file '%s'", importSpec.Path.Value, file)
				}
			}
		}

		dst.Decls = append(dst.Decls, src.Decls...)
		for _, importSpec := range missingImports {
			astutil.AddImport(fileSet, dst, strings.Replace(importSpec.Path.Value, "\"", "", -1))
		}
	}

	return &Generator{files: fileSet, code: dst}, nil
}

// Creates new framework source from byte array
func FromBytes(src []byte) (*Generator, error) {

	fileSet := token.NewFileSet()
	dst, err := parser.ParseFile(fileSet, "", src, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse framework code")
	}

	return &Generator{files: fileSet, code: dst}, nil
}

// Returns to the framework source code
func (gen *Generator) Bytes() ([]byte, error) {

	var buffer bytes.Buffer
	if err := format.Node(&buffer, gen.files, gen.code); err != nil {
		return nil, errors.Wrap(err, "failed to merge framework code")
	}

	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to format framework code")
	}

	return code, nil
}

// Appends source code
func (gen *Generator) Append(code []byte) error {

	src, err := parser.ParseFile(gen.files, "", code, 0)
	if err != nil {
		return errors.Wrap(err, "failed to parse appending code")
	}

	imports := astutil.Imports(gen.files, src)
	for _, importSpecs := range imports {
		for _, importSpec := range importSpecs {
			if ok := astutil.DeleteImport(gen.files, src, strings.Replace(importSpec.Path.Value, "\"", "", -1)); !ok {
				return errors.Errorf("failed to delete import %s", importSpec.Path.Value)
			}
		}
	}

	gen.code.Decls = append(gen.code.Decls, src.Decls...)
	return nil
}

// Renames package name
func (gen *Generator) RenamePackage(name string) {

	gen.code.Name = &ast.Ident{Name: name}
}

// Renames specific nodes for a given ast tree
func (gen *Generator) RenameTypes(pkgs []string, name string) error {

	rewriteFunc := func(n ast.Node) (ast.Node, bool) {
		switch x := n.(type) {
		case *ast.TypeSpec:
			x.Name.Name = replace(pkgs, x.Name.Name, name)
			return x, true
		case *ast.StructType:
			for i := 0; i < len(x.Fields.List); i++ {
				selectorExpr, ok := x.Fields.List[i].Type.(*ast.SelectorExpr)
				if ok {
					ident := selectorExpr.X.(*ast.Ident)
					if contains(pkgs, ident.Name) {
						x.Fields.List[i].Type = &ast.Ident{Name: selectorExpr.Sel.Name}
					}
				}
			}
			return x, true
		case *ast.SelectorExpr:
			ident, ok := x.X.(*ast.Ident)
			if ok {
				if contains(pkgs, ident.Name) {
					return &ast.Ident{Name: x.Sel.Name}, true
				}
			}
		}
		return n, true
	}

	astrewrite.Walk(gen.code, rewriteFunc)
	return nil
}

// Make certain identifiers private
func (gen *Generator) MakePrivate(identifiers []string) error {

	rewriteFunc := func(n ast.Node) (ast.Node, bool) {
		switch x := n.(type) {
		case *ast.Ident:
			if contains(identifiers, x.Name) {
				x.Name = stringutil.UnTitle(x.Name)
			}
			return x, true
		}
		return n, true
	}

	astrewrite.Walk(gen.code, rewriteFunc)
	return nil
}

func containsPath(pkgs []string, src string) bool {

	for _, pkg := range pkgs {
		if strings.Contains(src, pkg) {
			return true
		}
	}
	return false
}

func contains(pkgs []string, src string) bool {

	for _, pkg := range pkgs {
		if pkg == src {
			return true
		}
	}
	return false
}

func replace(pkgs []string, src string, dest string) string {

	for _, pkg := range pkgs {
		if pkg == src {
			return dest
		}
	}
	return src
}
