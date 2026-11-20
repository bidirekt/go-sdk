package reader

import (
	"cmp"
	"fmt"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/bidirekt/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

type structSchema struct {
	schema     model.NamedSchema
	referenced bool
	declared   bool
}

// A struct of the project with a marked field is written even when nothing references it, unless it only
// appears embedded in other structs: then it is part of them.
func (this *reading) readProjectStructs() {
	for _, pkg := range this.loaded.packages {
		if !this.loaded.inProject[pkg.PkgPath] || pkg.Types == nil {
			continue
		}
		if len(pkg.Errors) > 0 {
			if packageHasBidirektTag(pkg) {
				this.reportGoErrors(pkg)
			}
			continue
		}
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			typeName, isTypeName := scope.Lookup(name).(*types.TypeName)
			if !isTypeName || typeName.IsAlias() {
				continue
			}
			named, isNamed := typeName.Type().(*types.Named)
			if !isNamed {
				continue
			}
			fields, isStruct := named.Underlying().(*types.Struct)
			if !isStruct {
				continue
			}
			if named.TypeParams().Len() > 0 {
				if this.structToFields(fields).marked {
					this.addErrorAt(typeName.Pos(), "cannot translate generic type %q: declare a struct without type parameters, with its fields marked with bidirekt:\"field\", and use it as the type", typeName.Name())
				}
				continue
			}
			if schema := this.namedToSchema(named); schema != nil {
				schema.declared = true
			}
		}
	}
}

func (this *reading) namedToSchema(named *types.Named) *structSchema {
	typeName := named.Origin().Obj()
	if schema, found := this.schemas[typeName]; found {
		return schema
	}
	fields := named.Underlying().(*types.Struct)
	if typeName.Pkg() != nil {
		if pkg, loaded := this.loaded.byPath[typeName.Pkg().Path()]; loaded && len(pkg.Errors) > 0 {
			this.reportGoErrors(pkg)
		}
	}
	if !this.structToFields(fields).marked {
		this.schemas[typeName] = nil
		return nil
	}
	schema := &structSchema{schema: model.NamedSchema{Name: this.schemaName(typeName), Location: this.locationAt(typeName.Pos())}}
	this.schemas[typeName] = schema
	schema.schema.Node = this.structToObject(fields, schema.schema.Name)
	return schema
}

// owner is the schema's name, empty for an anonymous struct.
func (this *reading) structToObject(fields *types.Struct, owner string) model.Node {
	object := model.Object{}
	names := map[string][]*types.Var{}
	for _, visible := range this.structToFields(fields).visible {
		field := visible.field
		if !field.tag.marked {
			continue
		}
		name := field.tag.options.name
		if name == "" {
			name = visible.name
		}
		names[name] = append(names[name], field.variable)
		object.Properties = append(object.Properties, model.Property{Name: name, Node: this.fieldToNode(field), Location: this.locationAt(field.variable.Pos())})
	}
	where := "this anonymous struct"
	if owner != "" {
		where = fmt.Sprintf("schema %q", owner)
	}
	for _, property := range object.Properties {
		for _, variable := range names[property.Name] {
			if len(names[property.Name]) > 1 {
				this.addErrorAt(variable.Pos(), "field %q is declared more than once in %s: give each field a distinct name", property.Name, where)
			}
		}
	}
	return model.Node{Shape: object}
}

// The schemas of the contract: those referenced, those derived from bodies, and those of the project not
// only embedded, followed through the names their nodes reference.
func (this *reading) schemasToContract() {
	byName := map[string][]*structSchema{}
	var all []*structSchema
	for _, schema := range this.schemas {
		if schema != nil {
			all = append(all, schema)
			byName[schema.schema.Name] = append(byName[schema.schema.Name], schema)
		}
	}
	slices.SortFunc(all, func(a, b *structSchema) int {
		return cmp.Or(strings.Compare(a.schema.Location.File, b.schema.Location.File), cmp.Compare(a.schema.Location.Line, b.schema.Location.Line))
	})
	included := map[*structSchema]bool{}
	var pending []*structSchema
	for typeName, schema := range this.schemas {
		if schema != nil && (schema.referenced || (schema.declared && !this.embeddedNames[typeName])) {
			pending = append(pending, schema)
		}
	}
	for _, derived := range this.derivedSchemas {
		for _, name := range nodeToReferences(derived.Node) {
			pending = append(pending, byName[name]...)
		}
	}
	for len(pending) > 0 {
		schema := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if included[schema] {
			continue
		}
		included[schema] = true
		for _, name := range nodeToReferences(schema.schema.Node) {
			pending = append(pending, byName[name]...)
		}
	}
	for _, schema := range all {
		if included[schema] {
			this.result.Contract.Schemas = append(this.result.Contract.Schemas, schema.schema)
		}
	}
	this.result.Contract.Schemas = append(this.result.Contract.Schemas, this.derivedSchemas...)
}

func nodeToReferences(node model.Node) []string {
	switch shape := node.Shape.(type) {
	case model.Ref:
		return []string{shape.Schema}
	case model.Object:
		var names []string
		for _, property := range shape.Properties {
			names = append(names, nodeToReferences(property.Node)...)
		}
		return names
	case model.Array:
		return nodeToReferences(shape.Items)
	case model.Map:
		return nodeToReferences(shape.Values)
	case model.Union:
		var names []string
		for _, member := range shape.Members {
			names = append(names, nodeToReferences(member.Node)...)
		}
		return names
	}
	return nil
}

func packageHasBidirektTag(pkg *packages.Package) bool {
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		fields, isStruct := scope.Lookup(name).Type().Underlying().(*types.Struct)
		if !isStruct {
			continue
		}
		for index := range fields.NumFields() {
			if strings.Contains(fields.Tag(index), tagKey) {
				return true
			}
		}
	}
	return false
}

func (this *reading) locationAt(position token.Pos) model.SourceLocation {
	place := this.loaded.fset.Position(position)
	return model.SourceLocation{File: place.Filename, Line: place.Line}
}
