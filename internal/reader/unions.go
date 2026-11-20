package reader

import (
	"go/token"
	"go/types"

	"github.com/bidirekt/go-sdk/internal/model"
)

type unionField struct {
	written  types.Type
	methods  *types.Interface
	union    model.Union
	position token.Pos
	pkg      *types.Package
}

func reachesInterface(written types.Type) bool {
	switch shape := written.Underlying().(type) {
	case *types.Interface:
		return true
	case *types.Pointer:
		return reachesInterface(shape.Elem())
	case *types.Slice:
		return reachesInterface(shape.Elem())
	case *types.Array:
		return reachesInterface(shape.Elem())
	case *types.Map:
		return reachesInterface(shape.Elem())
	}
	return false
}

// A member of union= is a schema by its name in the contract: a struct that implements the field's
// interface, by value or by pointer. A name no schema declares is left to the check of every reference.
func (this *reading) readUnionMembers() {
	structs := map[string]types.Type{}
	others := map[string]bool{}
	for typeName, schema := range this.schemas {
		if schema == nil {
			continue
		}
		if _, isStruct := typeName.Type().Underlying().(*types.Struct); isStruct {
			structs[schema.schema.Name] = typeName.Type()
			continue
		}
		others[schema.schema.Name] = true
	}
	for _, name := range this.enumNames {
		others[name] = true
	}
	for _, derived := range this.derivedSchemas {
		others[derived.Name] = true
	}
	for _, field := range this.unions {
		for _, member := range field.union.Members {
			name := member.Node.Shape.(model.Ref).Schema
			memberType, isStruct := structs[name]
			switch {
			case isStruct && !types.Implements(memberType, field.methods) && !types.Implements(types.NewPointer(memberType), field.methods):
				this.addErrorAt(field.position, "the union= lists %q, which does not implement %q: add the methods of %q to %q or to its pointer", name, typeToText(field.written, field.pkg), typeToText(field.written, field.pkg), typeToText(memberType, field.pkg))
			case !isStruct && others[name]:
				this.addErrorAt(field.position, "the union= lists %q, which is not a struct: list structs with fields marked with bidirekt:\"field\"", name)
			}
		}
	}
}
