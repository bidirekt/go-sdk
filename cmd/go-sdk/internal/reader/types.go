package reader

import (
	"go/token"
	"go/types"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
)

type typePlace struct {
	position    token.Pos
	pkg         *types.Package
	union       *model.Union
	body        bool
	hasGoErrors bool
}

func (this *reading) fieldToNode(field *directField) model.Node {
	if cached, found := this.fieldNodes[field.variable]; found {
		return cached
	}
	written := field.variable.Type()
	options := field.tag.options
	var node model.Node
	switch {
	case len(field.tag.errors) > 0:
	case options.declared != nil:
		node = *options.declared
		node.Nullable = isPointer(written)
	case field.quoted:
		node = model.Node{Shape: model.Primitive{Type: "string"}, Nullable: isPointer(written)}
	case options.union != nil && !reachesInterface(written, map[types.Type]bool{}):
		this.addErrorAt(field.variable.Pos(), "cannot place the union= in %q: write it on a field whose type, items or values are an interface", typeToText(written, field.variable.Pkg()))
	default:
		place := typePlace{position: field.variable.Pos(), pkg: field.variable.Pkg(), union: options.union, hasGoErrors: this.reportGoErrorsOf(field.variable.Pkg())}
		node, _ = this.typeToNode(written, place)
	}
	if field.omitEmpty || field.omitZero {
		node.Nullable = false
	}
	node.Optional = field.omitZero || (field.omitEmpty && isOmittedWhenEmpty(written))
	if options.optional != nil {
		node.Optional = *options.optional
	}
	if options.nullable != nil {
		node.Nullable = *options.nullable
	}
	node.Description = options.description
	this.fieldNodes[field.variable] = node
	return node
}

func (this *reading) typeToNode(original types.Type, place typePlace) (model.Node, bool) {
	if this.reportGoErrorsOfType(original) {
		place.hasGoErrors = true
	}
	written := types.Unalias(original)
	// A struct is a ref to its schema; any other named type is written in place, so one that contains
	// itself would never end.
	if named, isNamed := written.(*types.Named); isNamed {
		if _, isStruct := named.Underlying().(*types.Struct); !isStruct {
			instance := types.TypeString(named, nil)
			if this.inProgress[instance] {
				this.addErrorAt(place.position, "cannot translate recursive type %q: only a struct can contain itself, so declare a struct with its fields marked with bidirekt:\"field\" and use it instead", typeToText(original, place.pkg))
				return model.Node{}, false
			}
			this.inProgress[instance] = true
			defer delete(this.inProgress, instance)
		}
	}
	if pointer, isPointer := written.Underlying().(*types.Pointer); isPointer {
		node, readable := this.typeToNode(pointer.Elem(), place)
		node.Nullable = true
		return node, readable
	}
	if named, isNamed := written.(*types.Named); isNamed {
		if isJSONNumber(named) {
			return primitiveNode("number"), true
		}
		if enum, found := this.enumNames[named.Obj()]; found {
			this.schemas[named.Obj()].referenced = true
			return model.Node{Shape: model.Ref{Schema: enum}}, true
		}
	}
	if hasMarshaler(written) {
		return this.untranslatable(original, place)
	}
	switch shape := written.Underlying().(type) {
	case *types.Basic:
		return this.basicToNode(original, shape, place)
	case *types.Slice:
		if isBytes(shape) {
			return this.untranslatable(original, place)
		}
		items, readable := this.typeToNode(shape.Elem(), place)
		return model.Node{Shape: model.Array{Items: items}}, readable
	case *types.Array:
		items, readable := this.typeToNode(shape.Elem(), place)
		return model.Node{Shape: model.Array{Items: items}}, readable
	case *types.Map:
		if basic, isBasic := shape.Key().Underlying().(*types.Basic); !isBasic || basic.Info()&types.IsString == 0 {
			if place.body {
				this.addErrorAt(place.position, "cannot translate %q in a body: a map key must be a string, so give the map a string key type", typeToText(original, place.pkg))
				return model.Node{}, false
			}
			this.addErrorAt(place.position, "cannot translate %q to a contract type: a map key must be a string, or declare it with bidirekt:\"field type=map values=<type>\"", typeToText(original, place.pkg))
			return model.Node{}, false
		}
		values, readable := this.typeToNode(shape.Elem(), place)
		return model.Node{Shape: model.Map{Values: values}}, readable
	case *types.Struct:
		return this.structTypeToNode(written, shape, place)
	case *types.Interface:
		if place.union != nil {
			this.unions = append(this.unions, unionField{written: original, methods: shape, union: *place.union, position: place.position, pkg: place.pkg})
			return model.Node{Shape: *place.union}, true
		}
	}
	return this.untranslatable(original, place)
}

func (this *reading) basicToNode(written types.Type, basic *types.Basic, place typePlace) (model.Node, bool) {
	info := basic.Info()
	switch {
	case basic.Kind() == types.Invalid:
		if !place.hasGoErrors {
			this.addErrorAt(place.position, "cannot translate %q: Go cannot resolve its type, so run go build and fix the errors it reports", typeToText(written, place.pkg))
		}
		return model.Node{}, false
	case info&types.IsBoolean != 0:
		return primitiveNode("boolean"), true
	case info&types.IsString != 0:
		return primitiveNode("string"), true
	case info&types.IsInteger != 0:
		return primitiveNode("integer"), true
	case info&types.IsFloat != 0:
		return primitiveNode("number"), true
	}
	return this.untranslatable(written, place)
}

func (this *reading) structTypeToNode(written types.Type, fields *types.Struct, place typePlace) (model.Node, bool) {
	named, isNamed := written.(*types.Named)
	if !isNamed {
		if !this.structToFields(fields).marked {
			this.addErrorAt(place.position, "cannot translate an anonymous struct without marked fields: mark the fields of the contract with bidirekt:\"field\"")
			return model.Node{}, false
		}
		return this.structToObject(fields, ""), true
	}
	if named.TypeArgs().Len() > 0 {
		this.addErrorAt(place.position, "cannot translate generic type %q: declare a struct without type parameters, with its fields marked with bidirekt:\"field\", and use it as the type", typeToText(written, place.pkg))
		return model.Node{}, false
	}
	schema := this.namedToSchema(named)
	if schema == nil {
		typeName := named.Obj()
		if typeName.Pkg() != nil && !this.loaded.local[typeName.Pkg().Path()] {
			this.addErrorAt(place.position, "cannot translate struct %q of package %q: declare a struct of your own with the fields of the contract marked with bidirekt:\"field\"", typeName.Name(), typeName.Pkg().Path())
			return model.Node{}, false
		}
		this.addErrorAt(place.position, "cannot translate struct %q: mark the fields of the contract with bidirekt:\"field\"", typeName.Name())
		return model.Node{}, false
	}
	schema.referenced = true
	return model.Node{Shape: model.Ref{Schema: schema.schema.Name}}, true
}

func (this *reading) reportGoErrorsOfType(written types.Type) (hasErrors bool) {
	for {
		switch typed := written.(type) {
		case *types.Alias:
			hasErrors = this.reportGoErrorsOf(typed.Obj().Pkg()) || hasErrors
			written = typed.Rhs()
		case *types.Named:
			return this.reportGoErrorsOf(typed.Obj().Pkg()) || hasErrors
		default:
			return hasErrors
		}
	}
}

func (this *reading) untranslatable(written types.Type, place typePlace) (model.Node, bool) {
	if place.body {
		this.addErrorAt(place.position, "cannot translate %q in a body: %sdeclare a struct with the fields of the contract marked with bidirekt:\"field\" and use it in the body", typeToText(written, place.pkg), marshalerToReason(written))
		return model.Node{}, false
	}
	this.addErrorAt(place.position, "cannot translate %q to a contract type: declare it with bidirekt:\"field type=<type>\"", typeToText(written, place.pkg))
	return model.Node{}, false
}

func hasMarshaler(written types.Type) bool {
	return typeToMarshaler(written) != ""
}

func typeToMarshaler(written types.Type) string {
	written = types.Unalias(written)
	for _, method := range []string{"MarshalJSON", "MarshalText"} {
		for _, candidate := range []types.Type{written, types.NewPointer(written)} {
			if types.NewMethodSet(candidate).Lookup(nil, method) != nil {
				return method
			}
		}
	}
	return ""
}

func marshalerToReason(written types.Type) string {
	if method := typeToMarshaler(written); method != "" {
		return "encoding/json writes it with its " + method + ", so "
	}
	return ""
}

func isJSONNumber(named *types.Named) bool {
	typeName := named.Obj()
	return typeName.Pkg() != nil && typeName.Pkg().Path() == "encoding/json" && typeName.Name() == "Number"
}

func isPointer(written types.Type) bool {
	_, pointer := written.Underlying().(*types.Pointer)
	return pointer
}

func isOmittedWhenEmpty(written types.Type) bool {
	switch written.Underlying().(type) {
	case *types.Basic, *types.Slice, *types.Map, *types.Pointer, *types.Interface:
		return true
	}
	return false
}

func primitiveNode(name string) model.Node {
	return model.Node{Shape: model.Primitive{Type: name}}
}

func typeToText(written types.Type, from *types.Package) string {
	return types.TypeString(written, func(pkg *types.Package) string {
		if pkg == from {
			return ""
		}
		return pkg.Name()
	})
}
