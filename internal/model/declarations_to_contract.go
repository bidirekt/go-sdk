package model

import (
	"fmt"
	"reflect"
	"slices"
)

func DeclarationsToContract(declarations Contract) (Contract, []Diagnostic) {
	schemas, errors := declarationsToSchemas(declarations.Schemas)
	provides, providedErrors := providedToMerged(declarations.Provides)
	consumes, consumedErrors := consumedToMerged(declarations.Consumes)
	contract := Contract{Provides: provides, Consumes: consumes, Schemas: schemas}
	errors = slices.Concat(errors, providedErrors, consumedErrors, contractToMissingReferences(contract, declarations.Schemas))
	return contract, errors
}

type reference struct {
	schema   string
	location SourceLocation
}

func contractToMissingReferences(contract Contract, declarations []NamedSchema) []Diagnostic {
	declared := map[string]bool{}
	for _, declaration := range declarations {
		declared[declaration.Name] = true
	}
	var references []reference
	for _, resource := range append(slices.Clone(contract.Provides), contract.Consumes...) {
		bodies := []Body{}
		if resource.Request != nil {
			bodies = append(bodies, *resource.Request)
		}
		for _, response := range resource.Responses {
			bodies = append(bodies, response.Body)
		}
		for _, body := range bodies {
			if body.Schema != "" {
				references = append(references, reference{schema: body.Schema, location: body.Location})
			}
		}
	}
	for _, schema := range contract.Schemas {
		references = append(references, nodeToReferences(schema.Node, schema.Location)...)
	}
	var errors []Diagnostic
	for _, each := range references {
		if !declared[each.schema] {
			errors = append(errors, Diagnostic{
				Message:  fmt.Sprintf("schema %q is not declared: reference a schema by its name in the contract", each.schema),
				Location: each.location,
			})
		}
	}
	return errors
}

func nodeToReferences(node Node, location SourceLocation) []reference {
	switch shape := node.Shape.(type) {
	case Ref:
		return []reference{{schema: shape.Schema, location: location}}
	case Object:
		var references []reference
		for _, property := range shape.Properties {
			references = append(references, nodeToReferences(property.Node, property.Location)...)
		}
		return references
	case Array:
		return nodeToReferences(shape.Items, location)
	case Map:
		return nodeToReferences(shape.Values, location)
	case Union:
		var references []reference
		for _, member := range shape.Members {
			references = append(references, nodeToReferences(member.Node, location)...)
		}
		return references
	}
	return nil
}

func declarationsToSchemas(declarations []NamedSchema) ([]NamedSchema, []Diagnostic) {
	var names []string
	groups := map[string][]NamedSchema{}
	for _, declaration := range declarations {
		if _, found := groups[declaration.Name]; !found {
			names = append(names, declaration.Name)
		}
		groups[declaration.Name] = append(groups[declaration.Name], declaration)
	}
	var schemas []NamedSchema
	var errors []Diagnostic
	for _, name := range names {
		group := groups[name]
		if isOneDerivedSchema(group) {
			schemas = append(schemas, group[0])
			continue
		}
		message := collisionToMessage(name, group)
		for _, declaration := range group {
			errors = append(errors, Diagnostic{Message: message, Location: declaration.Location})
		}
	}
	return schemas, errors
}

func isOneDerivedSchema(group []NamedSchema) bool {
	for _, declaration := range group {
		if len(group) > 1 && (!declaration.Derived || !reflect.DeepEqual(nodeToShape(declaration.Node), nodeToShape(group[0].Node))) {
			return false
		}
	}
	return true
}

func collisionToMessage(name string, group []NamedSchema) string {
	lists, derived := 0, 0
	for _, declaration := range group {
		if !declaration.Derived {
			continue
		}
		derived++
		if _, isList := declaration.Node.Shape.(Array); isList {
			lists++
		}
	}
	switch {
	case derived == 0:
		return fmt.Sprintf("schema name %q is declared more than once: rename one with bidirekt.Schema(..., bidirekt.Named(\"<name>\"))", name)
	case lists == len(group):
		return fmt.Sprintf("lists named %q differ: give each list a distinct name", name)
	case lists > 0:
		return fmt.Sprintf("list name %q is also the name of another schema: give the list a distinct name", name)
	}
	return fmt.Sprintf("schema name %q is declared more than once: give each body a distinct name with bidirekt.Named", name)
}

// nodeToShape drops where each property is declared, so two bodies written alike in two places are the same.
func nodeToShape(node Node) Node {
	switch shape := node.Shape.(type) {
	case Object:
		properties := make([]Property, 0, len(shape.Properties))
		for _, property := range shape.Properties {
			properties = append(properties, Property{Name: property.Name, Node: nodeToShape(property.Node)})
		}
		node.Shape = Object{Properties: properties}
	case Array:
		node.Shape = Array{Items: nodeToShape(shape.Items)}
	case Map:
		node.Shape = Map{Values: nodeToShape(shape.Values)}
	case Union:
		members := make([]UnionMember, 0, len(shape.Members))
		for _, member := range shape.Members {
			members = append(members, UnionMember{Key: member.Key, Node: nodeToShape(member.Node)})
		}
		node.Shape = Union{Members: members}
	}
	return node
}
