package reader

import (
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func (this *reading) exprToConstant(pkg *packages.Package, expr ast.Expr) (constant.Value, bool) {
	value := pkg.TypesInfo.Types[expr].Value
	if value == nil {
		this.addError(expr, "%q cannot be read without running the code: write a literal or a constant", types.ExprString(expr))
		return nil, false
	}
	return value, true
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
