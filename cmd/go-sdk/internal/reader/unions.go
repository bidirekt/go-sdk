package reader

import (
	"go/token"
	"go/types"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
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

// A member of union= is a schema by its name in the contract: a struct of the contract that implements the
// field's interface, by value or by pointer. A name no schema declares is left to the check of every
// reference.
func (this *reading) readUnionMembers(contractTypeNames []*types.TypeName) {
	known := map[string]bool{}
	structs := map[string][]types.Type{}
	for _, typeName := range contractTypeNames {
		name := this.schemas[typeName].schema.Name
		known[name] = true
		if _, isStruct := typeName.Type().Underlying().(*types.Struct); isStruct {
			structs[name] = append(structs[name], typeName.Type())
		}
	}
	for _, derived := range this.derivedSchemas {
		known[derived.Name] = true
	}
	for _, field := range this.unions {
		for _, member := range field.union.Members {
			name := member.Node.Shape.(model.Ref).Schema
			memberTypes, isStruct := structs[name]
			if !isStruct && known[name] {
				this.addErrorAt(field.position, "the union= lists %q, which is not a struct: list structs with fields marked with bidirekt:\"field\"", name)
			}
			for _, memberType := range memberTypes {
				if !types.Implements(memberType, field.methods) && !types.Implements(types.NewPointer(memberType), field.methods) {
					this.addErrorAt(field.position, "the union= lists %q, which does not implement %q: add the methods of %q to %q or to its pointer", name, typeToText(field.written, field.pkg), typeToText(field.written, field.pkg), typeToText(memberType, field.pkg))
				}
			}
		}
	}
}
