package reader

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func (this *reading) packageToCallsWithArgumentsInPlace(pkg *packages.Package) []sdkCall {
	var calls []sdkCall
	for _, call := range packageToCalls(pkg) {
		if isMultiValueCall(pkg.TypesInfo, call.call) {
			this.addError(call.call, "arguments of bidirekt.%s must be written in place, one by one: write them instead of %q", call.function, types.ExprString(call.call.Args[0]))
			continue
		}
		calls = append(calls, call)
	}
	return calls
}

func (this *reading) readOptionsInPlace(each packageCalls) {
	info := each.pkg.TypesInfo
	inPlace := map[*ast.CallExpr]bool{}
	for _, call := range each.calls {
		options := []string{"Named"}
		if call.function == "Provides" || call.function == "Consumes" {
			options = []string{"Request", "Response"}
		}
		for _, arg := range call.call.Args {
			if option, _, found := callToSDKOption(info, arg, options...); found {
				inPlace[option] = true
			}
		}
	}
	for _, call := range each.calls {
		if inPlace[call.call] {
			continue
		}
		switch call.function {
		case "Request", "Response":
			this.addError(call.call, "bidirekt.%s must be written in place, as an argument of bidirekt.Provides or bidirekt.Consumes", call.function)
		case "Named":
			this.addError(call.call, "bidirekt.Named must be written in place, as an argument of bidirekt.Request, bidirekt.Response or bidirekt.Schema")
		}
	}
}

// optionsToName reads the one bidirekt.Named an owner may take: named is false without it.
func (this *reading) optionsToName(pkg *packages.Package, owner sdkCall, options []ast.Expr) (name string, named bool, readable bool) {
	if owner.call.Ellipsis.IsValid() {
		this.addError(owner.call, "options of bidirekt.%s must be written in place, one by one: write them without ...", owner.function)
		return "", false, false
	}
	readable = true
	for _, option := range options {
		call, _, found := callToSDKOption(pkg.TypesInfo, option, "Named")
		if !found {
			this.addError(option, "%q is not an option written in place: write bidirekt.Named(\"<name>\") as an argument of bidirekt.%s", types.ExprString(option), owner.function)
			readable = false
			continue
		}
		value, isString := this.exprToString(pkg, call.Args[0])
		switch {
		case !isString:
			readable = false
		case named:
			this.addError(call, "bidirekt.Named is declared more than once in this bidirekt.%s: keep one", owner.function)
			readable = false
		case value == "":
			this.addError(call, "bidirekt.Named is empty: write the name of the schema")
			readable = false
		}
		name, named = value, true
	}
	return name, named, readable
}
