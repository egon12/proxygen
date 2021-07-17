package proxygen

import (
	"fmt"
	"go/ast"
)

type (
	transformer struct{}
)

func (t *transformer) Transform(i *InterfaceType, suffix string) (Proxy, error) {
	receiver := Var{"t", "*" + i.Name + suffix}

	funcs, err := t.transformFunctions(i.Ast, receiver, i.Name)
	if err != nil {
		return Proxy{}, err
	}

	return Proxy{
		PackageName: i.PackageName,
		BaseType:    i.Name,
		Type:        i.Name + suffix,
		Funcs:       funcs,
	}, nil
}

func (t *transformer) transformFunctions(it *ast.InterfaceType, receiver Var, baseType string) ([]Func, error) {
	var err error
	funcs := make([]Func, len(it.Methods.List))

	for i, m := range it.Methods.List {
		ft, ok := m.Type.(*ast.FuncType)
		if !ok {
			return nil, fmt.Errorf("casting method(%T) to functype failed", m)
		}

		funcs[i].Name = m.Names[0].Name
		funcs[i].Receiver = receiver
		funcs[i].BaseType = baseType
		funcs[i].Params, funcs[i].Return, err = t.getMultiVar(ft)
		if err != nil {
			return nil, err
		}
		funcs[i].FixEmptyParams()
	}

	return funcs, nil
}

func (t *transformer) getMultiVar(f *ast.FuncType) (params, returns MultiVar, err error) {
	params, err = t.transformFieldList(f.Params)
	if err != nil {
		err = fmt.Errorf("transform params failed: %w", err)
		return
	}

	returns, err = t.transformFieldList(f.Results)
	if err != nil {
		err = fmt.Errorf("transform returns failed: %w", err)
		return
	}

	return
}

func (t *transformer) transformFieldList(fields *ast.FieldList) (MultiVar, error) {
	result := make([]Var, len(fields.List))

	for i, f := range fields.List {
		name := ""
		if len(f.Names) > 0 {
			name = f.Names[0].Name
		}

		typeName, err := t.transformFieldType(f.Type)
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

func (t *transformer) transformFieldType(ex ast.Expr) (string, error) {
	switch ft := ex.(type) {
	case *ast.Ident:
		return ft.Name, nil
	case *ast.SelectorExpr:
		res, err := t.transformFieldType(ft.X)
		return res + "." + ft.Sel.Name, err
	case *ast.StarExpr:
		res, err := t.transformFieldType(ft.X)
		return "*" + res, err
	case *ast.ArrayType:
		res, err := t.transformFieldType(ft.Elt)
		return "[]" + res, err
	case *ast.MapType:
		key, err := t.transformFieldType(ft.Key)
		if err != nil {
			return "", err
		}
		value, err := t.transformFieldType(ft.Value)
		if err != nil {
			return "", err
		}
		return "map[" + key + "]" + value, nil

	default:
		return "", fmt.Errorf("cannot transform fieldtype %T", t)
	}
}
