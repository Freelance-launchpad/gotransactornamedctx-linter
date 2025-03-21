package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	singlechecker.Main(Analyzer)
}

var Analyzer = &analysis.Analyzer{
	Name:     "gotransactornamedctx",
	Doc:      "Checks that transactor.WithinTransaction has a named context parameter.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		callExpr := node.(*ast.CallExpr)

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if selectorExpr.Sel == nil || selectorExpr.Sel.Name != "WithinTransaction" {
			return
		}

		selectorIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		// TODO: check if this change when renaming the import
		if selectorIdent.Name != "transactor" {
			return
		}

		if len(callExpr.Args) != 2 {
			return
		}

		funcLit, ok := callExpr.Args[1].(*ast.FuncLit)
		if !ok {
			return
		}

		if funcLit.Type == nil || funcLit.Type.Params == nil || len(funcLit.Type.Params.List) == 0 {
			return
		}

		field := funcLit.Type.Params.List[0]
		if field == nil {
			return
		}

		if len(field.Names) == 0 {
			pass.Reportf(node.Pos(), "transactor function has unnamed context parameter\n")
		}
		return
	})

	return nil, nil
}
