package model

import (
	"fmt"
	"slices"
)

func providedToMerged(provides []Resource) ([]Resource, []Diagnostic) {
	var resources []Resource
	var errors []Diagnostic
	for _, group := range resourcesToGroups(provides, func(resource Resource) [3]string {
		return [3]string{"", resource.Method, resource.Path}
	}) {
		resources = append(resources, group[0])
		if len(group) == 1 {
			continue
		}
		message := fmt.Sprintf(`"%s %s" is provided by more than one bidirekt.Provides: declare each method and path in one bidirekt.Provides`, group[0].Method, group[0].Path)
		for _, resource := range group {
			errors = append(errors, Diagnostic{Message: message, Location: resource.Location})
		}
	}
	return resources, errors
}

// The bidirekt.Consumes of the same provider, method and path are one resource with the statuses of all.
func consumedToMerged(consumes []Resource) ([]Resource, []Diagnostic) {
	var resources []Resource
	var errors []Diagnostic
	for _, group := range resourcesToGroups(consumes, func(resource Resource) [3]string {
		return [3]string{resource.Provider, resource.Method, resource.Path}
	}) {
		first := group[0]
		named := fmt.Sprintf(`"%s %s" consumed from "%s"`, first.Method, first.Path, first.Provider)
		for _, other := range group[1:] {
			errors = append(errors, requestsToConflicts(first, other, named)...)
		}
		merged := first
		merged.Responses = nil
		for _, declaration := range group {
			for _, response := range declaration.Responses {
				index := slices.IndexFunc(merged.Responses, func(kept Response) bool { return kept.Status == response.Status })
				if index == -1 {
					merged.Responses = append(merged.Responses, response)
					continue
				}
				kept := merged.Responses[index]
				if kept.Body.Schema == response.Body.Schema {
					continue
				}
				message := fmt.Sprintf(
					"response %d of %s is declared as %s and %s: declare one schema with every field the service reads",
					response.Status, named, schemaToText(kept.Body.Schema), schemaToText(response.Body.Schema),
				)
				errors = append(errors, Diagnostic{Message: message, Location: kept.Body.Location}, Diagnostic{Message: message, Location: response.Body.Location})
			}
		}
		resources = append(resources, merged)
	}
	return resources, errors
}

func requestsToConflicts(first, other Resource, named string) []Diagnostic {
	switch {
	case first.Request == nil && other.Request == nil:
		return nil
	case first.Request == nil || other.Request == nil:
		message := fmt.Sprintf("request of %s is declared in one bidirekt.Consumes and not in another: declare the same request in every bidirekt.Consumes", named)
		return []Diagnostic{{Message: message, Location: requestToLocation(first)}, {Message: message, Location: requestToLocation(other)}}
	case first.Request.Schema == other.Request.Schema:
		return nil
	}
	message := fmt.Sprintf(
		"request of %s is declared as %s and %s: declare the same request in every bidirekt.Consumes",
		named, schemaToText(first.Request.Schema), schemaToText(other.Request.Schema),
	)
	return []Diagnostic{{Message: message, Location: first.Request.Location}, {Message: message, Location: other.Request.Location}}
}

func requestToLocation(resource Resource) SourceLocation {
	if resource.Request == nil {
		return resource.Location
	}
	return resource.Request.Location
}

func schemaToText(schema string) string {
	if schema == "" {
		return "null"
	}
	return `"` + schema + `"`
}

func resourcesToGroups(resources []Resource, resourceToKey func(Resource) [3]string) [][]Resource {
	var groups [][]Resource
	indexByKey := map[[3]string]int{}
	for _, resource := range resources {
		key := resourceToKey(resource)
		index, found := indexByKey[key]
		if !found {
			index = len(groups)
			indexByKey[key] = index
			groups = append(groups, nil)
		}
		groups[index] = append(groups[index], resource)
	}
	return groups
}
