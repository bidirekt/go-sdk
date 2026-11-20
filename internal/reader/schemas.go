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
		for _, named := range scopeToStructs(pkg.Types.Scope()) {
			marked := this.structToFields(named.Underlying().(*types.Struct)).marked
			if named.TypeParams().Len() > 0 {
				if marked {
					this.addErrorAt(named.Obj().Pos(), "cannot translate generic type %q: declare a struct without type parameters, with its fields marked with bidirekt:\"field\", and use it as the type", named.Obj().Name())
				}
				continue
			}
			if method := typeToMarshaler(named); method != "" {
				if marked {
					this.addErrorAt(named.Obj().Pos(), "struct %q has its own %s, so encoding/json never writes its marked fields: remove bidirekt:\"field\" from its fields, or remove the method", named.Obj().Name(), method)
				}
				continue
			}
			if schema := this.namedToSchema(named); schema != nil {
				schema.declared = true
			}
		}
	}
}

func scopeToStructs(scope *types.Scope) []*types.Named {
	var structs []*types.Named
	for _, name := range scope.Names() {
		typeName, isTypeName := scope.Lookup(name).(*types.TypeName)
		if !isTypeName || typeName.IsAlias() {
			continue
		}
		named, isNamed := typeName.Type().(*types.Named)
		if !isNamed {
			continue
		}
		if _, isStruct := named.Underlying().(*types.Struct); isStruct {
			structs = append(structs, named)
		}
	}
	return structs
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
	outer := this.inProgress
	this.inProgress = map[string]bool{}
	schema.schema.Node = this.structToObject(fields, schema.schema.Name)
	this.inProgress = outer
	return schema
}

// A schema named only in a tag (type=, items=, values=, union=) is a named struct of the loaded packages or
// of the packages they import, found by its name in the contract; a struct found this way may name more.
func (this *reading) readTagReferences() {
	var candidates map[string][]*types.Named
	for found := true; found; {
		found = false
		known := this.contractSchemaNames()
		for _, name := range this.referencedSchemaNames() {
			if known[name] {
				continue
			}
			if candidates == nil {
				candidates = this.packagesToStructsByName()
			}
			for _, named := range candidates[name] {
				if schema := this.namedToSchema(named); schema != nil {
					schema.referenced = true
					found = true
				}
			}
			delete(candidates, name)
		}
	}
}

func (this *reading) packagesToStructsByName() map[string][]*types.Named {
	byName := map[string][]*types.Named{}
	visited := map[*types.Package]bool{}
	var pending []*types.Package
	for _, pkg := range this.loaded.packages {
		if pkg.Types != nil {
			pending = append(pending, pkg.Types)
		}
	}
	for len(pending) > 0 {
		pkg := pending[0]
		pending = pending[1:]
		if visited[pkg] {
			continue
		}
		visited[pkg] = true
		pending = append(pending, pkg.Imports()...)
		for _, named := range scopeToStructs(pkg.Scope()) {
			if named.TypeParams().Len() == 0 {
				name := this.schemaName(named.Obj())
				byName[name] = append(byName[name], named)
			}
		}
	}
	return byName
}

func (this *reading) contractSchemaNames() map[string]bool {
	names := map[string]bool{}
	for _, schema := range this.schemas {
		if schema != nil {
			names[schema.schema.Name] = true
		}
	}
	for _, derived := range this.derivedSchemas {
		names[derived.Name] = true
	}
	return names
}

func (this *reading) referencedSchemaNames() []string {
	var nodes []model.NamedSchema
	for _, schema := range this.schemas {
		if schema != nil {
			nodes = append(nodes, schema.schema)
		}
	}
	var names []string
	for _, schema := range append(nodes, this.derivedSchemas...) {
		for _, reference := range model.NodeToReferences(schema.Node, schema.Location) {
			names = append(names, reference.Schema)
		}
	}
	slices.Sort(names)
	return slices.Compact(names)
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
		node := this.fieldToNode(field)
		if visible.throughPointer && field.tag.options.optional == nil {
			node.Optional = true
		}
		object.Properties = append(object.Properties, model.Property{Name: name, Node: node, Location: this.locationAt(field.variable.Pos())})
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
		for _, reference := range model.NodeToReferences(derived.Node, derived.Location) {
			pending = append(pending, byName[reference.Schema]...)
		}
	}
	for len(pending) > 0 {
		schema := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if included[schema] {
			continue
		}
		included[schema] = true
		for _, reference := range model.NodeToReferences(schema.schema.Node, schema.schema.Location) {
			pending = append(pending, byName[reference.Schema]...)
		}
	}
	for _, schema := range all {
		if included[schema] {
			this.result.Contract.Schemas = append(this.result.Contract.Schemas, schema.schema)
		}
	}
	this.result.Contract.Schemas = append(this.result.Contract.Schemas, this.derivedSchemas...)
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
