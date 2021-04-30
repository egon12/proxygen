package proxygen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strconv"
)

type (
	Collector struct {
		file *ast.File
	}

	InterfaceTransformer struct {
		name        string
		packageName string
		source      *ast.InterfaceType
	}

	InterfaceTransformConfig struct {
		PackageName string
	}
)

var (
	Default = InterfaceTransformConfig{}
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

func (c *Collector) FindInterface(name string) (*InterfaceTransformer, error) {
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
				return &InterfaceTransformer{name, c.file.Name.Name, it}, nil
			}
		}
	}

	return nil, fmt.Errorf("cannot find interface %s", name)
}

func (t *InterfaceTransformer) Transform() (Proxy, error) {
	receiverType := "*" + t.name + "Tracer"
	receiver := Var{"t", receiverType}

	funcs, _ := t.transformFunctions(t.source)

	for i := range funcs {
		funcs[i].Receiver = receiver

		for j := range funcs[i].Params {
			if funcs[i].Params[j].Name == "" {
				funcs[i].Params[j].Name = "arg" + strconv.Itoa(j)
			}
		}
	}

	// overwrite to remove the pointer (*)
	receiver.Type = t.name + "Tracer"

	return Proxy{
		PackageName:  t.packageName,
		Receiver:     receiver,
		OriginalType: t.name,
		Funcs:        funcs,
	}, nil
}

func (t *InterfaceTransformer) transformFunctions(it *ast.InterfaceType) ([]Func, error) {
	res := make([]Func, len(it.Methods.List))

	for i, m := range it.Methods.List {
		ft, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil, fmt.Errorf("casting method(%T) to functype failed", m)
		}

		// TODO Check error?
		funcName := m.Names[0].Name

		// TODO check error
		res[i], _ = t.transformFunction(funcName, ft)
	}

	return res, nil
}

func (c *InterfaceTransformer) transformFunction(name string, f *ast.FuncType) (Func, error) {
	return Func{
		Name:   name,
		Params: c.transformFieldList(f.Params),
		Return: c.transformFieldList(f.Results),
	}, nil
}

func (c *InterfaceTransformer) transformFieldList(fields *ast.FieldList) MultiVar {
	result := make([]Var, len(fields.List))

	for i, f := range fields.List {
		typeID, ok := f.Type.(*ast.Ident)
		if !ok {
			// TODO do something
		}

		name := ""
		if len(f.Names) > 0 {
			name = f.Names[0].Name
		}

		result[i] = Var{
			Name: name,
			Type: typeID.Name,
		}
	}

	return result
}
