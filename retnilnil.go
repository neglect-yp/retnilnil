package retnilnil

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "retnilnil"
	doc  = "Retnilnil is a static analysis tool to detect `return nil, nil`"
)

var errorType = types.Universe.Lookup("error").Type()

var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

type context struct {
	pass        *analysis.Pass
	commentMaps *comment.Maps
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	commentMaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	nodes := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Nodes(nodes, func(n ast.Node, push bool) bool {
		if !push {
			return false
		}

		if file, ok := n.(*ast.File); ok && strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			return false
		}

		ctx := &context{
			pass:        pass,
			commentMaps: &commentMaps,
		}

		body, ok := getBodyIfTarget(ctx, n)
		if !ok {
			return true
		}

		walk(ctx, body)

		return true
	})

	return nil, nil
}

func getBodyIfTarget(ctx *context, n ast.Node) (*ast.BlockStmt, bool) {
	switch n := n.(type) {
	case *ast.FuncDecl:
		return n.Body, isResultsTypeMatched(ctx, n.Type.Results) && !hasCommentAboutNilNil(ctx, n)
	case *ast.FuncLit:
		return n.Body, isResultsTypeMatched(ctx, n.Type.Results)
	}

	return nil, false
}

func isResultsTypeMatched(ctx *context, results *ast.FieldList) (ok bool) {
	if results.NumFields() != 2 {
		return false
	}

	t1 := ctx.pass.TypesInfo.TypeOf(results.List[0].Type)
	if _, ok := t1.(*types.Pointer); !ok && !types.IsInterface(t1) {
		return false
	}

	t2 := ctx.pass.TypesInfo.TypeOf(results.List[1].Type)
	if t2 != errorType {
		return false
	}

	return true
}

func hasCommentAboutNilNil(ctx *context, decl *ast.FuncDecl) (ok bool) {
	for _, comment := range ctx.commentMaps.Comments(decl) {
		if strings.Contains(comment.Text(), "nil, nil") {
			return true
		}
	}

	return false
}

func walk(ctx *context, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.ReturnStmt:
		reportIfDetected(ctx, stmt)
	case *ast.IfStmt:
		walk(ctx, stmt.Body)
		walk(ctx, stmt.Else)
	case *ast.ForStmt:
		walk(ctx, stmt.Body)
	case *ast.RangeStmt:
		walk(ctx, stmt.Body)
	case *ast.SwitchStmt:
		walk(ctx, stmt.Body)
	case *ast.TypeSwitchStmt:
		walk(ctx, stmt.Body)
	case *ast.CaseClause:
		for _, stmt := range stmt.Body {
			walk(ctx, stmt)
		}
	case *ast.BlockStmt:
		for _, stmt := range stmt.List {
			walk(ctx, stmt)
		}
	}
}

func reportIfDetected(ctx *context, stmt *ast.ReturnStmt) {
	if ctx.commentMaps.Ignore(stmt, name) {
		return
	}

	if len(stmt.Results) != 2 {
		return
	}

	v1, ok := stmt.Results[0].(*ast.Ident)
	if !ok {
		return
	}

	v2, ok := stmt.Results[1].(*ast.Ident)
	if !ok {
		return
	}

	if v1.Name == "nil" && v2.Name == "nil" {
		ctx.pass.Reportf(stmt.Pos(), "`return nil, nil` should be avoided. Please consider to use a pointer to a zero value or an appropriate error like ErrNotFound")
	}
}
