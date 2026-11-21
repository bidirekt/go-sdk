package reader

import (
	"go/ast"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

type sdkCall struct {
	function string
	call     *ast.CallExpr
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
				if function, found := callToSDKFunction(pkg.TypesInfo, call); found {
					calls = append(calls, sdkCall{function: function, call: call})
				}
			}
			return true
		})
	}
	return calls
}

func callToSDKFunction(info *types.Info, call *ast.CallExpr) (string, bool) {
	function := typeutil.StaticCallee(info, call)
	if function == nil || function.Pkg() == nil || function.Pkg().Path() != sdkPath {
		return "", false
	}
	return function.Name(), true
}

func callToSDKOption(info *types.Info, expr ast.Expr, functions ...string) (*ast.CallExpr, string, bool) {
	call, isCall := ast.Unparen(expr).(*ast.CallExpr)
	if !isCall {
		return nil, "", false
	}
	function, found := callToSDKFunction(info, call)
	if !found || !slices.Contains(functions, function) {
		return nil, "", false
	}
	return call, function, true
}

func isMultiValueCall(info *types.Info, call *ast.CallExpr) bool {
	if len(call.Args) != 1 {
		return false
	}
	_, isTuple := info.TypeOf(call.Args[0]).(*types.Tuple)
	return isTuple
}
