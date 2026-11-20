package reader

import (
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func (this *reading) exprToConstant(pkg *packages.Package, expr ast.Expr) (constant.Value, bool) {
	hasGoErrors := this.reportGoErrorsOfUses(pkg, expr)
	value := pkg.TypesInfo.Types[expr].Value
	switch {
	case value == nil:
		this.addError(expr, "%q cannot be read without running the code: write a literal or a constant", types.ExprString(expr))
		return nil, false
	case value.Kind() == constant.Unknown:
		if !hasGoErrors {
			this.addError(expr, "%q has no value Go can compute: run go build and fix the errors it reports", types.ExprString(expr))
		}
		return nil, false
	}
	return value, true
}

func (this *reading) reportGoErrorsOfUses(pkg *packages.Package, expr ast.Expr) (hasErrors bool) {
	ast.Inspect(expr, func(node ast.Node) bool {
		if ident, isIdent := node.(*ast.Ident); isIdent {
			if object := pkg.TypesInfo.Uses[ident]; object != nil && this.reportGoErrorsOf(object.Pkg()) {
				hasErrors = true
			}
		}
		return true
	})
	return hasErrors
}

func (this *reading) exprToString(pkg *packages.Package, expr ast.Expr) (string, bool) {
	value, isConstant := this.exprToConstant(pkg, expr)
	if !isConstant || value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(value), true
}

func (this *reading) exprToInt(pkg *packages.Package, expr ast.Expr) (int, bool) {
	value, isConstant := this.exprToConstant(pkg, expr)
	if !isConstant {
		return 0, false
	}
	number, exact := constant.Int64Val(constant.ToInt(value))
	return int(number), exact
}
