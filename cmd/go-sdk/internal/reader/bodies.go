package reader

import (
	"go/ast"
	"go/types"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

type bodyPlace struct {
	expr  ast.Expr
	typed typePlace
	name  string
	named bool
}

func (this *reading) readBody(pkg *packages.Package, option sdkCall, expr ast.Expr, options []ast.Expr) (model.Body, bool) {
	body := model.Body{Location: this.location(expr)}
	name, named, readable := this.optionsToName(pkg, option, options)
	if !readable {
		return body, false
	}
	typed := pkg.TypesInfo.Types[expr]
	if typed.IsNil() {
		if named {
			this.addError(expr, "a nil body has no schema to name: remove bidirekt.Named")
			return body, false
		}
		return body, true
	}
	place := bodyPlace{expr: expr, typed: typePlace{position: expr.Pos(), pkg: pkg.Types, body: true}, name: name, named: named}
	body.Schema, readable = this.typeToBodySchema(typed.Type, place)
	return body, readable
}

// The body is the static type of the expression, a pointer counted by the type it points to.
func (this *reading) typeToBodySchema(original types.Type, place bodyPlace) (string, bool) {
	pointed := deref(original)
	if this.reportGoErrorsOfType(pointed) {
		place.typed.hasGoErrors = true
	}
	written := types.Unalias(pointed)
	if named, isNamed := written.(*types.Named); isNamed && !hasMarshaler(written) {
		switch named.Underlying().(type) {
		case *types.Struct:
			return this.namedToBodySchema(written, place)
		case *types.Slice, *types.Array, *types.Map:
			if !isBytes(named.Underlying()) {
				return this.collectionToBodySchema(named, place)
			}
		}
	}
	if hasMarshaler(written) {
		return this.unreadableBody(pointed, place)
	}
	switch shape := written.(type) {
	case *types.Slice:
		if !isBytes(shape) {
			return this.listToBodySchema(shape.Elem(), place)
		}
	case *types.Array:
		return this.listToBodySchema(shape.Elem(), place)
	case *types.Map:
		if _, isInterface := shape.Elem().Underlying().(*types.Interface); !isInterface {
			return this.derivedToBodySchema(written, "a map", place)
		}
	case *types.Struct:
		return this.derivedToBodySchema(written, "an anonymous struct", place)
	}
	return this.unreadableBody(pointed, place)
}

func (this *reading) namedToBodySchema(written types.Type, place bodyPlace) (string, bool) {
	node, readable := this.typeToNode(written, place.typed)
	if !readable {
		return "", false
	}
	schema := node.Shape.(model.Ref).Schema
	if place.named {
		this.addError(place.expr, "the body is the schema %q, whose name bidirekt.Named cannot change: rename it with bidirekt.Schema(..., bidirekt.Named(%q))", schema, place.name)
		return "", false
	}
	return schema, true
}

func (this *reading) collectionToBodySchema(named *types.Named, place bodyPlace) (string, bool) {
	if named.TypeArgs().Len() > 0 {
		this.addError(place.expr, "cannot translate generic type %q: declare a slice or map type without type parameters, as type Pets []Pet, and use it as the body", typeToText(named, place.typed.pkg))
		return "", false
	}
	typeName := named.Obj()
	if place.named {
		this.addError(place.expr, "the body is the schema %q, whose name bidirekt.Named cannot change: rename it with bidirekt.Schema(..., bidirekt.Named(%q))", this.schemaName(typeName), place.name)
		return "", false
	}
	if schema, found := this.schemas[typeName]; found {
		schema.referenced = true
		return schema.schema.Name, true
	}
	node, readable := this.typeToNode(named.Underlying(), place.typed)
	if !readable {
		return "", false
	}
	schema := &structSchema{
		schema:     model.NamedSchema{Name: this.schemaName(typeName), Node: node, Location: this.locationAt(typeName.Pos())},
		referenced: true,
	}
	this.schemas[typeName] = schema
	return schema.schema.Name, true
}

func (this *reading) listToBodySchema(items types.Type, place bodyPlace) (string, bool) {
	node, readable := this.typeToNode(items, place.typed)
	if !readable {
		return "", false
	}
	name := place.name
	if !place.named {
		ref, isRef := node.Shape.(model.Ref)
		if !isRef {
			this.addError(place.expr, "the list body %q has items without a schema name to name it after: name it with bidirekt.Named(\"<name>\")", types.ExprString(place.expr))
			return "", false
		}
		name = ref.Schema + "List"
	}
	this.addDerivedSchema(name, model.Node{Shape: model.Array{Items: node}}, place)
	return name, true
}

func (this *reading) derivedToBodySchema(written types.Type, kind string, place bodyPlace) (string, bool) {
	if !place.named {
		this.addError(place.expr, "%s body has no schema name: name it with bidirekt.Named(\"<name>\")", kind)
		return "", false
	}
	node, readable := this.typeToNode(written, place.typed)
	if !readable {
		return "", false
	}
	this.addDerivedSchema(place.name, node, place)
	return place.name, true
}

func (this *reading) addDerivedSchema(name string, node model.Node, place bodyPlace) {
	this.derivedSchemas = append(this.derivedSchemas, model.NamedSchema{Name: name, Node: node, Location: this.location(place.expr), Derived: true})
}

func (this *reading) unreadableBody(written types.Type, place bodyPlace) (string, bool) {
	this.addError(place.expr, "cannot translate a body of type %q: %sdeclare a struct with the fields of the contract marked with bidirekt:\"field\" and use it as the body", typeToText(written, place.typed.pkg), marshalerToReason(written))
	return "", false
}

func isBytes(collection types.Type) bool {
	slice, isSlice := collection.Underlying().(*types.Slice)
	if !isSlice {
		return false
	}
	basic, isBasic := slice.Elem().Underlying().(*types.Basic)
	return isBasic && basic.Kind() == types.Byte
}
