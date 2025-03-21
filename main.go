package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "--" {
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	if selectorExpr.Sel == nil || selectorExpr.Sel.Name != "WithinTransaction" {
		return v
	}

	selectorIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return v
	}

	// TODO: check if this change when renaming the import
	if selectorIdent.Name != "transactor" {
		return v
	}

	if len(callExpr.Args) != 2 {
		return v
	}

	funcLit, ok := callExpr.Args[1].(*ast.FuncLit)
	if !ok {
		return v
	}

	if funcLit.Type == nil || funcLit.Type.Params == nil || len(funcLit.Type.Params.List) == 0 {
		return v
	}

	field := funcLit.Type.Params.List[0]
	if field == nil {
		return v
	}

	if len(field.Names) == 0 {
		fmt.Printf("%s: transactor function has unnamed context parameter\n", v.fset.Position(node.Pos()))
	}

	return v
}
