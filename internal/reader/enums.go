package reader

import (
	"go/ast"
	"go/types"
	"slices"

	"github.com/bidirekt/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

type schemaDeclaration struct {
	typeName  *types.TypeName
	call      *ast.CallExpr
	inProject bool
}

func (this *reading) readSchemaDeclaration(pkg *packages.Package, declaration sdkCall) {
	name, named, readable := this.optionsToName(pkg, declaration, declaration.call.Args[1:])
	if !readable {
		return
	}
	value := declaration.call.Args[0]
	written := types.Unalias(deref(pkg.TypesInfo.Types[value].Type))
	if !named {
		this.addError(declaration.call, "bidirekt.Schema needs a name: write bidirekt.Schema(%s, bidirekt.Named(\"<name>\"))", types.ExprString(value))
		return
	}
	typed, isNamed := written.(*types.Named)
	if !isNamed {
		this.addNotASchema(declaration.call, written, pkg.Types)
		return
	}
	typeName := typed.Obj()
	for _, earlier := range this.schemaDeclarations {
		if earlier.typeName == typeName {
			message := "bidirekt.Schema is declared more than once for %q: keep one"
			this.addError(earlier.call, message, typeName.Name())
			this.addError(declaration.call, message, typeName.Name())
			return
		}
	}
	this.schemaNames[typeName] = name
	this.schemaDeclarations = append(this.schemaDeclarations, schemaDeclaration{typeName: typeName, call: declaration.call, inProject: this.loaded.inProject[pkg.PkgPath]})
}

// A bidirekt.Schema names a struct with marked fields, a named list or map used as a body, or a type of
// bidirekt.Enum. Outside the project, a list or map may be a body only of another service.
func (this *reading) readSchemaDeclarationTypes() {
	for _, declaration := range this.schemaDeclarations {
		typed := declaration.typeName.Type().(*types.Named)
		if _, isEnum := this.enumNames[declaration.typeName]; isEnum {
			continue
		}
		switch shape := typed.Underlying().(type) {
		case *types.Slice, *types.Array, *types.Map:
			if this.schemas[declaration.typeName] != nil || !declaration.inProject {
				continue
			}
		case *types.Struct:
			if typed.TypeParams().Len() == 0 && this.structToFields(shape).marked {
				continue
			}
		}
		this.addNotASchema(declaration.call, typed, declaration.typeName.Pkg())
	}
}

func (this *reading) addNotASchema(call *ast.CallExpr, written types.Type, from *types.Package) {
	this.addError(call, "bidirekt.Schema is on %q, which is not a schema nor an enum: put it on a struct with fields marked with bidirekt:\"field\" or on a type declared with bidirekt.Enum", typeToText(written, from))
}

func (this *reading) readEnum(pkg *packages.Package, declaration sdkCall) {
	if declaration.call.Ellipsis.IsValid() {
		this.addError(declaration.call, "values of bidirekt.Enum must be written in place, one by one: write them without ...")
		return
	}
	signature := pkg.TypesInfo.TypeOf(declaration.call.Fun).(*types.Signature)
	valueType := signature.Params().At(0).Type().(*types.Slice).Elem()
	named, isNamed := types.Unalias(valueType).(*types.Named)
	if !isNamed {
		this.addError(declaration.call, "values of bidirekt.Enum have no named type: declare a string type, as type PetStatus string, and give the values that type")
		return
	}
	typeName := named.Obj()
	if len(declaration.call.Args) == 0 {
		this.addError(declaration.call, "bidirekt.Enum of %q has no values: write the values of the enum as its arguments", typeName.Name())
		return
	}
	var members []string
	readable := true
	for _, arg := range declaration.call.Args {
		member, isString := this.exprToString(pkg, arg)
		switch {
		case !isString:
			readable = false
		case member == "":
			this.addError(arg, "a value of bidirekt.Enum is empty: write the value or remove it")
			readable = false
		case slices.Contains(members, member):
			this.addError(arg, "value %q is declared more than once in bidirekt.Enum: keep one", member)
			readable = false
		}
		members = append(members, member)
	}
	if earlier, declared := this.enumDeclarations[typeName]; declared {
		message := "bidirekt.Enum is declared more than once for %q: keep one"
		this.addError(earlier, message, typeName.Name())
		this.addError(declaration.call, message, typeName.Name())
		return
	}
	this.enumDeclarations[typeName] = declaration.call
	if !readable {
		return
	}
	name := this.schemaName(typeName)
	this.enumNames[typeName] = name
	this.schemas[typeName] = &structSchema{
		schema:   model.NamedSchema{Name: name, Node: model.Node{Shape: model.Enum{Members: members}}, Location: this.location(declaration.call)},
		declared: this.loaded.inProject[pkg.PkgPath],
	}
}
