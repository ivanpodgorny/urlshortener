package osexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"regexp"
)

// Analyzer анализатор, который проверяет, что в функции main пакета main
// не используется прямой вызов os.Exit.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit calls in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	r := regexp.MustCompile("^// Code generated .* DO NOT EDIT.$")
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}
		if len(file.Comments) > 0 && r.MatchString(file.Comments[0].List[0].Text) {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				if n.Name.Name != "main" {
					return false
				}
			case *ast.ExprStmt:
				if sel, ok := n.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" && sel.Sel.Name == "Exit" {
						pass.Reportf(n.Pos(), "call of os.Exit in main")

						return false
					}
				}
			}

			return true
		})
	}

	return nil, nil
}
