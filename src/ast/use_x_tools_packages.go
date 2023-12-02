package ast

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// DetectCallErrorsIs は errors.Is の呼び出し箇所を検出します
func DetectCallErrorsIs(dir string) ([]token.Position, error) {
	var detectPositions []token.Position
	conf := &packages.Config{
		Mode: packages.NeedTypes |
			packages.NeedSyntax,
		Dir: dir,
	}
	pkgs, err := packages.Load(conf, "./...")
	if err != nil {
		return nil, fmt.Errorf("ast: load packages: %w", err)
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// 関数呼び出しを探す
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				var expression, identName string
				switch fun := callExpr.Fun.(type) {
				case *ast.Ident: // Is() など `.` を使わない呼び出し
					identName = fun.Name
				case *ast.SelectorExpr: // errors.Is() など `.` を使った呼び出し
					identName = fun.Sel.Name
					if ident, ok := fun.X.(*ast.Ident); ok {
						expression = ident.Name
					}
				}
				if identName != "Is" {
					return true
				}
				// import spec と合わせて errors パッケージの Is 呼び出しなのか検証
				for _, importSpec := range file.Imports {
					if importSpec.Path.Value != `"errors"` {
						continue
					}
					switch expression {
					case "errors": // import "errors"
						if importSpec.Name != nil {
							return true
						}
					case "": // import . "errors"
						if importSpec.Name == nil || importSpec.Name.Name != "." {
							return true
						}
					default: // import alias "errors"
						if importSpec.Name == nil || importSpec.Name.Name != expression {
							return true
						}
					}
					pos := pkg.Fset.Position(n.Pos())
					detectPositions = append(detectPositions, pos)
				}
				return true
			})
		}
	}
	return detectPositions, nil
}
