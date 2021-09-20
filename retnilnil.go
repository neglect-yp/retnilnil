package retnilnil

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "Retnilnil is a static analysis tool to detect `return nil, nil`"

var errorType = types.Universe.Lookup("error").Type()

var Analyzer = &analysis.Analyzer{
	Name: "retnilnil",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodes := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodes, func(n ast.Node) {
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		if !isSignatureMatched(pass, decl) {
			return
		}

		for _, stmt := range decl.Body.List {
			if stmt, ok := stmt.(*ast.ReturnStmt); ok {
				v1, ok := stmt.Results[0].(*ast.Ident)
				if !ok {
					continue
				}

				v2, ok := stmt.Results[1].(*ast.Ident)
				if !ok {
					continue
				}

				if v1.Name == "nil" && v2.Name == "nil" {
					pass.Reportf(stmt.Pos(), "return nil, nil")
				}
			}
		}
	})

	return nil, nil
}

func isSignatureMatched(pass *analysis.Pass, decl *ast.FuncDecl) (ok bool) {
	results := decl.Type.Results

	if results.NumFields() != 2 {
		return false
	}

	t1 := pass.TypesInfo.TypeOf(results.List[0].Type)
	if _, ok := t1.(*types.Pointer); !ok && !types.IsInterface(t1) {
		return false
	}

	t2 := pass.TypesInfo.TypeOf(results.List[1].Type)
	if t2 != errorType {
		return false
	}

	return true
}
