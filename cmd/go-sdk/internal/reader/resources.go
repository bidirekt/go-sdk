package reader

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

func (this *reading) readResource(pkg *packages.Package, declaration sdkCall) {
	if declaration.call.Ellipsis.IsValid() {
		this.addError(declaration.call, "options of bidirekt.%s must be written in place, one by one: write them without ...", declaration.function)
		return
	}
	fixed := 2
	if declaration.function == "Consumes" {
		fixed = 3
	}
	values := make([]string, 0, fixed)
	for _, arg := range declaration.call.Args[:fixed] {
		value, isString := this.exprToString(pkg, arg)
		if !isString {
			return
		}
		values = append(values, value)
	}
	resource := model.Resource{Location: this.location(declaration.call)}
	if declaration.function == "Consumes" {
		resource.Provider, values = values[0], values[1:]
	}
	resource.Method, resource.Path = values[0], values[1]
	readable := true
	if upper := strings.ToUpper(resource.Method); upper != resource.Method {
		this.addError(declaration.call.Args[fixed-2], "method %q must be uppercase: write %q", resource.Method, upper)
		readable = false
	}
	for _, arg := range declaration.call.Args[fixed:] {
		readable = this.readResourceOption(pkg, declaration, arg, &resource) && readable
	}
	if !readable {
		return
	}
	if declaration.function == "Consumes" {
		this.result.Contract.Consumes = append(this.result.Contract.Consumes, resource)
		return
	}
	this.result.Contract.Provides = append(this.result.Contract.Provides, resource)
}

func (this *reading) readResourceOption(pkg *packages.Package, declaration sdkCall, arg ast.Expr, resource *model.Resource) bool {
	call, function, found := callToSDKOption(pkg.TypesInfo, arg, "Request", "Response")
	if !found {
		this.addError(arg, "%q is not an option written in place: write bidirekt.Request(...) or bidirekt.Response(...) as an argument of bidirekt.%s", types.ExprString(arg), declaration.function)
		return false
	}
	if isMultiValueCall(pkg.TypesInfo, call) {
		return false
	}
	option := sdkCall{function: function, call: call}
	if function == "Request" {
		if resource.Request != nil {
			this.addError(call, "bidirekt.%s declares more than one bidirekt.Request: keep one", declaration.function)
			return false
		}
		body, readable := this.readBody(pkg, option, call.Args[0], call.Args[1:])
		resource.Request = &body
		return readable
	}
	status, statusReadable := this.readStatus(pkg, call.Args[0])
	body, bodyReadable := this.readBody(pkg, option, call.Args[1], call.Args[2:])
	if statusReadable && slices.ContainsFunc(resource.Responses, func(response model.Response) bool { return response.Status == status }) {
		this.addError(call, "response status %d is declared more than once in this bidirekt.%s: keep one", status, declaration.function)
		return false
	}
	resource.Responses = append(resource.Responses, model.Response{Status: status, Body: body})
	return statusReadable && bodyReadable
}

func (this *reading) readStatus(pkg *packages.Package, expr ast.Expr) (int, bool) {
	if literal, isLiteral := ast.Unparen(expr).(*ast.BasicLit); isLiteral && literal.Kind == token.INT && isLegacyOctal(literal.Value) {
		this.addError(literal, "response status %q must be written as %s", literal.Value, legacyOctalToDecimal(literal.Value))
		return 0, false
	}
	return this.exprToInt(pkg, expr)
}

func isLegacyOctal(literal string) bool {
	return len(literal) > 1 && literal[0] == '0' && strings.Trim(literal, "0123456789_") == ""
}

func legacyOctalToDecimal(literal string) string {
	digits := strings.TrimLeft(strings.ReplaceAll(literal, "_", ""), "0")
	if digits == "" {
		return "0"
	}
	return digits
}
