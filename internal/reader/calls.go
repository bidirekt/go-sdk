package reader

import (
	"go/ast"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

type sdkCall struct {
	function string
	call     *ast.CallExpr
	ident    *ast.Ident
}

func packageToCalls(pkg *packages.Package) []sdkCall {
	files := slices.Clone(pkg.Syntax)
	slices.SortFunc(files, func(a, b *ast.File) int {
		return strings.Compare(pkg.Fset.Position(a.Pos()).Filename, pkg.Fset.Position(b.Pos()).Filename)
	})
	var calls []sdkCall
	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			if call, isCall := node.(*ast.CallExpr); isCall {
				if function, ident, found := callToSDKFunction(pkg.TypesInfo, call); found {
					calls = append(calls, sdkCall{function: function, call: call, ident: ident})
				}
			}
			return true
		})
	}
	return calls
}

func callToSDKFunction(info *types.Info, call *ast.CallExpr) (string, *ast.Ident, bool) {
	fun := ast.Unparen(call.Fun)
	switch index := fun.(type) {
	case *ast.IndexExpr:
		fun = ast.Unparen(index.X)
	case *ast.IndexListExpr:
		fun = ast.Unparen(index.X)
	}
	var ident *ast.Ident
	switch name := fun.(type) {
	case *ast.SelectorExpr:
		ident = name.Sel
	case *ast.Ident:
		ident = name
	default:
		return "", nil, false
	}
	function, isFunction := info.Uses[ident].(*types.Func)
	if !isFunction || function.Pkg() == nil || function.Pkg().Path() != sdkPath {
		return "", nil, false
	}
	return function.Name(), ident, true
}

func callToSDKOption(info *types.Info, expr ast.Expr, functions ...string) (*ast.CallExpr, string, bool) {
	call, isCall := ast.Unparen(expr).(*ast.CallExpr)
	if !isCall {
		return nil, "", false
	}
	function, _, found := callToSDKFunction(info, call)
	if !found || !slices.Contains(functions, function) {
		return nil, "", false
	}
	return call, function, true
}
