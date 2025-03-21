package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(Analyzer)
}

var Analyzer = &analysis.Analyzer{
	Name: "gotransactornamedctx",
	Doc:  "Checks that transactor.WithinTransaction has a named context parameter.",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if selectorExpr.Sel == nil || selectorExpr.Sel.Name != "WithinTransaction" {
			return true
		}

		selectorIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			return true
		}

		// TODO: check if this change when renaming the import
		if selectorIdent.Name != "transactor" {
			return true
		}

		if len(callExpr.Args) != 2 {
			return true
		}

		funcLit, ok := callExpr.Args[1].(*ast.FuncLit)
		if !ok {
			return true
		}

		if funcLit.Type == nil || funcLit.Type.Params == nil || len(funcLit.Type.Params.List) == 0 {
			return true
		}

		field := funcLit.Type.Params.List[0]
		if field == nil {
			return true
		}

		if len(field.Names) == 0 {
			pass.Reportf(node.Pos(), "transactor function has unnamed context parameter\n")
		}
		return true
	}
	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
