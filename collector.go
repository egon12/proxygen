package proxygen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

type (
	// Collector is the first entry door of this logic.
	// To start generated code we need to instantiate
	// this collector with NewCollector. Then fill the
	// filename that this Collector need to read
	Collector struct {
		file *ast.File
	}

	InterfaceType struct {
		Name        string
		PackageName string
		Ast         *ast.InterfaceType
	}
)

// NewCollector will return new Collector
func NewCollector() *Collector {
	return &Collector{}
}

// Load will make the program read the file, then get ast
// from it
func (c *Collector) Load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot open %s: %v", filename, err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, string(b), parser.AllErrors)
	if err != nil {
		return fmt.Errorf("parse error %s: %v", filename, err)
	}

	c.file = f

	return nil
}

// FindInterface will return InterfaceType that can be transformed by transformer
func (c *Collector) FindInterface(name string) (*InterfaceType, error) {
	for _, decl := range c.file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, s := range genDecl.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if ts.Name.Name == name {
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("%s is not an interface it's a %v", name, ts.Type)
				}
				return &InterfaceType{
					Name:        name,
					PackageName: c.file.Name.Name,
					Ast:         it,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find interface %s", name)
}
