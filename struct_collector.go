package proxygen

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type (
	StructMethods struct {
		StructName  string
		PackageName string
		Ast         []*ast.FuncDecl
	}
)

func FindStructMethods(path string, structName string) (StructMethods, error) {
	result := StructMethods{StructName: structName}

	// TODO add filter test_file
	pkgs, err := parser.ParseDir(token.NewFileSet(), path, nil, parser.AllErrors)
	if err != nil {
		return result, err
	}

	for _, pkg := range pkgs {
		ast.Inspect(pkg, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.Package, *ast.File:
				return true

			case *ast.FuncDecl:
				f := node.(*ast.FuncDecl)
				ok := checkFuncDecl(f, structName)
				if ok {
					result.Ast = append(result.Ast, f)
				}
			}
			return false
		})
	}

	return result, nil
}

func checkFuncDecl(f *ast.FuncDecl, structName string) bool {
	// only process that have receiver
	if f.Recv == nil {
		return false
	}

	// only process that have reciver identity
	if len(f.Recv.List) < 1 {
		return false
	}

	t := f.Recv.List[0].Type

	// only process method that use pointer to receiver
	star, ok := t.(*ast.StarExpr)
	if !ok {
		return false
	}

	// only process correct pointer to receiver
	recvrID, ok := star.X.(*ast.Ident)
	if !ok {
		return false
	}

	// only process to methods that have receiver
	// pointer to A
	if recvrID.Name != structName {
		return false
	}

	return true
}
