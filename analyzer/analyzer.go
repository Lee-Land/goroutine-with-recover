package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goroutinewithrecover",
	Doc:      "Checks that goroutine has recover in defer function",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectorC := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.GoStmt)(nil),
	}
	inspectorC.Preorder(nodeFilter, func(node ast.Node) {
		gostat := node.(*ast.GoStmt)
		var r bool
		switch gostat.Call.Fun.(type) {
		case *ast.FuncLit:
			funcLit := gostat.Call.Fun.(*ast.FuncLit)
			r = hasRecover(funcLit.Body)
		case *ast.Ident:
			id := gostat.Call.Fun.(*ast.Ident)
			fd, ok := id.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return
			}
			r = hasRecover(fd.Body)
		default:

		}
		if !r {
			pass.Reportf(node.Pos(), "goroutine should have recover in defer func")
		}
	})
	return nil, nil
}

func hasRecover(bs *ast.BlockStmt) bool {
	for _, blockStmt := range bs.List {
		deferStmt, ok := blockStmt.(*ast.DeferStmt)
		if !ok {
			return false
		}
		switch deferStmt.Call.Fun.(type) {
		case *ast.SelectorExpr:
			selectorExpr := deferStmt.Call.Fun.(*ast.SelectorExpr)
			if "Recover" == selectorExpr.Sel.Name {
				return true
			}
		case *ast.FuncLit:
			fl := deferStmt.Call.Fun.(*ast.FuncLit)
			for i := range fl.Body.List {
				stmt := fl.Body.List[i]
				switch stmt.(type) {
				case *ast.ExprStmt:
					exprStmt := stmt.(*ast.ExprStmt)
					if isRecoverExpr(exprStmt.X) {
						return true
					}
				case *ast.IfStmt:
					is := stmt.(*ast.IfStmt)
					as, ok := is.Init.(*ast.AssignStmt)
					if !ok {
						continue
					}
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}
				case *ast.AssignStmt:
					as := stmt.(*ast.AssignStmt)
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}
				}
			}
		}
	}
	return false
}

func isRecoverExpr(expr ast.Expr) bool {
	ac, ok := expr.(*ast.CallExpr) // r := recover()
	if !ok {
		return false
	}
	id, ok := ac.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if "recover" == id.Name {
		return true
	}
	return false
}
