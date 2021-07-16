package proxygen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

type (
	Collector struct {
		file *ast.File
	}

	interfaceType struct {
		Name        string
		PackageName string
		Ast         *ast.InterfaceType
	}
)

func NewCollector() *Collector {
	return &Collector{}
}

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

func (c *Collector) FindInterface(name string) (*interfaceType, error) {
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
				return &interfaceType{
					Name:        name,
					PackageName: c.file.Name.Name,
					Ast:         it,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find interface %s", name)
}
