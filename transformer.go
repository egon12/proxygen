package proxygen

import (
	"fmt"
	"go/ast"
	"strconv"
)

func (t *InterfaceTransformer) Transform() (Proxy, error) {
	receiver := Var{"t", "*" + t.name}

	funcs, err := t.transformFunctions(t.source)
	if err != nil {
		return Proxy{}, err
	}

	for i := range funcs {
		funcs[i].BaseType = t.name
		funcs[i].Receiver = receiver

		for j := range funcs[i].Params {
			if funcs[i].Params[j].Name == "" {
				funcs[i].Params[j].Name = "arg" + strconv.Itoa(j)
			}
		}
	}

	p := Proxy{
		PackageName: t.packageName,
		Receiver:    receiver,
		Funcs:       funcs,
		BaseType:    t.name,
	}

	p.SetRecieverTypeSuffix("Tracer")

	return p, nil
}

func (t *InterfaceTransformer) transformFunctions(it *ast.InterfaceType) ([]Func, error) {
	var err error
	res := make([]Func, len(it.Methods.List))

	for i, m := range it.Methods.List {
		ft, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil, fmt.Errorf("casting method(%T) to functype failed", m)
		}

		// TODO Check error?
		funcName := m.Names[0].Name

		// TODO check error
		res[i], err = t.transformFunction(funcName, ft)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (c *InterfaceTransformer) transformFunction(name string, f *ast.FuncType) (Func, error) {
	var err error
	result := Func{}

	result.Name = name

	result.Params, err = c.transformFieldList(f.Params)
	if err != nil {
		return result, err
	}

	result.Return, err = c.transformFieldList(f.Results)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *InterfaceTransformer) transformFieldList(fields *ast.FieldList) (MultiVar, error) {
	result := make([]Var, len(fields.List))

	for i, f := range fields.List {
		name := ""
		if len(f.Names) > 0 {
			name = f.Names[0].Name
		}

		typeName, err := c.transformFieldType(f.Type)
		if err != nil {
			return nil, fmt.Errorf("error transform field %s: %v", name, err)
		}

		result[i] = Var{
			Name: name,
			Type: typeName,
		}
	}

	return result, nil
}

func (c *InterfaceTransformer) transformFieldType(t ast.Expr) (string, error) {
	switch ft := t.(type) {
	case *ast.Ident:
		return ft.Name, nil
	case *ast.SelectorExpr:
		res, err := c.transformFieldType(ft.X)
		return res + "." + ft.Sel.Name, err
	case *ast.StarExpr:
		res, err := c.transformFieldType(ft.X)
		return "*" + res, err
	case *ast.ArrayType:
		res, err := c.transformFieldType(ft.Elt)
		return "[]" + res, err
	case *ast.MapType:
		key, err := c.transformFieldType(ft.Key)
		if err != nil {
			return "", err
		}
		value, err := c.transformFieldType(ft.Value)
		if err != nil {
			return "", err
		}
		return "map[" + key + "]" + value, nil

	default:
		return "", fmt.Errorf("cannot transform fieldtype %T", t)
	}
}
