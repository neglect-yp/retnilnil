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

		walk(pass, decl.Body)
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

func walk(pass *analysis.Pass, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.ReturnStmt:
		reportIfDetected(pass, stmt)
	case *ast.IfStmt:
		walk(pass, stmt.Body)
		walk(pass, stmt.Else)
	case *ast.ForStmt:
		walk(pass, stmt.Body)
	case *ast.RangeStmt:
		walk(pass, stmt.Body)
	case *ast.SwitchStmt:
		walk(pass, stmt.Body)
	case *ast.TypeSwitchStmt:
		walk(pass, stmt.Body)
	case *ast.CaseClause:
		for _, stmt := range stmt.Body {
			walk(pass, stmt)
		}
	case *ast.BlockStmt:
		for _, stmt := range stmt.List {
			walk(pass, stmt)
		}
	}
}

func reportIfDetected(pass *analysis.Pass, stmt *ast.ReturnStmt) {
	v1, ok := stmt.Results[0].(*ast.Ident)
	if !ok {
		return
	}

	v2, ok := stmt.Results[1].(*ast.Ident)
	if !ok {
		return
	}

	if v1.Name == "nil" && v2.Name == "nil" {
		pass.Reportf(stmt.Pos(), "return nil, nil")
	}
}
