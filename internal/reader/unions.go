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

func reachesInterface(written types.Type, visited map[types.Type]bool) bool {
	if visited[written] {
		return false
	}
	visited[written] = true
	switch shape := written.Underlying().(type) {
	case *types.Interface:
		return true
	case *types.Pointer:
		return reachesInterface(shape.Elem(), visited)
	case *types.Slice:
		return reachesInterface(shape.Elem(), visited)
	case *types.Array:
		return reachesInterface(shape.Elem(), visited)
	case *types.Map:
		return reachesInterface(shape.Elem(), visited)
	}
	return false
}

// A member of union= is a schema by its name in the contract: a struct that implements the field's
// interface, by value or by pointer. A name no schema declares is left to the check of every reference.
func (this *reading) readUnionMembers() {
	known := this.contractSchemaNames()
	structs := map[string]types.Type{}
	for typeName, schema := range this.schemas {
		if _, isStruct := typeName.Type().Underlying().(*types.Struct); isStruct && schema != nil {
			structs[schema.schema.Name] = typeName.Type()
		}
	}
	for _, field := range this.unions {
		for _, member := range field.union.Members {
			name := member.Node.Shape.(model.Ref).Schema
			memberType, isStruct := structs[name]
			switch {
			case isStruct && !types.Implements(memberType, field.methods) && !types.Implements(types.NewPointer(memberType), field.methods):
				this.addErrorAt(field.position, "the union= lists %q, which does not implement %q: add the methods of %q to %q or to its pointer", name, typeToText(field.written, field.pkg), typeToText(field.written, field.pkg), typeToText(memberType, field.pkg))
			case !isStruct && known[name]:
				this.addErrorAt(field.position, "the union= lists %q, which is not a struct: list structs with fields marked with bidirekt:\"field\"", name)
			}
		}
	}
}
