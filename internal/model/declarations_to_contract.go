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

type Reference struct {
	Schema   string
	Location SourceLocation
}

func contractToMissingReferences(contract Contract, declarations []NamedSchema) []Diagnostic {
	declared := map[string]bool{}
	for _, declaration := range declarations {
		declared[declaration.Name] = true
	}
	var references []Reference
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
				references = append(references, Reference(body))
			}
		}
	}
	for _, schema := range contract.Schemas {
		references = append(references, NodeToReferences(schema.Node, schema.Location)...)
	}
	var errors []Diagnostic
	for _, each := range references {
		if !declared[each.Schema] {
			errors = append(errors, Diagnostic{
				Message:  fmt.Sprintf("schema %q is not declared: reference a schema by its name in the contract", each.Schema),
				Location: each.Location,
			})
		}
	}
	return errors
}

func NodeToReferences(node Node, location SourceLocation) []Reference {
	switch shape := node.Shape.(type) {
	case Ref:
		return []Reference{{Schema: shape.Schema, Location: location}}
	case Object:
		var references []Reference
		for _, property := range shape.Properties {
			references = append(references, NodeToReferences(property.Node, property.Location)...)
		}
		return references
	case Array:
		return NodeToReferences(shape.Items, location)
	case Map:
		return NodeToReferences(shape.Values, location)
	case Union:
		var references []Reference
		for _, member := range shape.Members {
			references = append(references, NodeToReferences(member.Node, location)...)
		}
		return references
	}
	return nil
}

func declarationsToSchemas(declarations []NamedSchema) ([]NamedSchema, []Diagnostic) {
	var schemas []NamedSchema
	var errors []Diagnostic
	for _, group := range ValuesToGroups(declarations, func(declaration NamedSchema) string { return declaration.Name }) {
		if isOneDerivedSchema(group) {
			schemas = append(schemas, group[0])
			continue
		}
		message := collisionToMessage(group[0].Name, group)
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
