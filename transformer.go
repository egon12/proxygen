package proxygen

import (
	"fmt"
	"go/ast"
)

func (t *InterfaceTransformer) Transform() (Proxy, error) {
	receiver := Var{"t", "*" + t.name}

	funcs, err := t.transformFunctions(t.source, receiver, t.name)
	if err != nil {
		return Proxy{}, err
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

func (t *InterfaceTransformer) TransformCollectedInterface(ci CollectedInterface) (Proxy, error) {
	receiver := Var{"t", "*" + ci.Name}

	funcs, err := t.transformFunctions(t.source, receiver, ci.Name)
	if err != nil {
		return Proxy{}, err
	}

	p := Proxy{
		PackageName: ci.PackageName,
		BaseType:    ci.Name,
		Receiver:    receiver,
		Funcs:       funcs,
	}

	p.SetRecieverTypeSuffix("Tracer")

	return p, nil
}

func (t *InterfaceTransformer) transformFunctions(it *ast.InterfaceType, receiver Var, baseType string) ([]Func, error) {
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

func (c *InterfaceTransformer) getMultiVar(f *ast.FuncType) (params, returns MultiVar, err error) {
	params, err = c.transformFieldList(f.Params)
	if err != nil {
		err = fmt.Errorf("transform params failed: %w", err)
		return
	}

	returns, err = c.transformFieldList(f.Results)
	if err != nil {
		err = fmt.Errorf("transform returns failed: %w", err)
		return
	}

	return
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
