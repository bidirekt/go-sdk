package reader

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/bidirekt/go-sdk/internal/model"
)

const tagKey = "bidirekt"

var fieldKeys = []string{"name", "type", "items", "values", "optional", "nullable", "description", "union"}

var primitiveTypes = []string{"string", "integer", "float", "number", "boolean"}

var schemaNamePattern = regexp.MustCompile(`^[^\s<>\[\],:={}'"]+$`)

var unionMemberPattern = regexp.MustCompile(`^([^\s:,]+)\s*:\s*([^\s:,]+)$`)

var descriptionPattern = regexp.MustCompile(`^[A-Za-z0-9]+( [A-Za-z0-9]+)*$`)

var wordPattern = regexp.MustCompile(`^\w+$`)

var mapShortForm = regexp.MustCompile(`^map\[[^\]]*\](.+)$`)

type fieldOptions struct {
	name        string
	declared    *model.Node
	union       *model.Union
	optional    *bool
	nullable    *bool
	description string
}

type structTag struct {
	marked   bool
	options  fieldOptions
	errors   []string
	warnings []string
	jsonTag  string
}

// rawToStructTag reads the bidirekt and json keys as reflect.StructTag.Lookup does, and the raw tag where
// Lookup cannot read it: a malformed tag that names bidirekt is an error, a key close to bidirekt a warning.
func rawToStructTag(raw string) structTag {
	tag := structTag{}
	tag.jsonTag = reflect.StructTag(raw).Get("json")
	value, found := reflect.StructTag(raw).Lookup(tagKey)
	keys, wellFormed := rawToKeys(raw)
	for _, key := range keys {
		switch {
		case key == tagKey:
		case strings.Contains(key, tagKey) && !wordPattern.MatchString(key):
			wellFormed = false
		case strings.HasPrefix(strings.ToLower(key), "bidire"):
			tag.warnings = append(tag.warnings, fmt.Sprintf(`%q is not a bidirekt tag: write bidirekt:"field"`, key))
		}
	}
	if !found {
		if !wellFormed && strings.Contains(raw, tagKey) {
			tag.errors = append(tag.errors, fmt.Sprintf("cannot read the struct tag `%s`: write each key as key:\"value\", separated by spaces, as in `json:\"name\" bidirekt:\"field\"`", raw))
		}
		return tag
	}
	tag.marked = true
	options, err := valueToFieldOptions(value)
	if err != "" {
		tag.errors = append(tag.errors, err)
	}
	tag.options = options
	return tag
}

func rawToKeys(raw string) ([]string, bool) {
	var keys []string
	for {
		raw = strings.TrimLeft(raw, " ")
		if raw == "" {
			return keys, true
		}
		end := 0
		for end < len(raw) && raw[end] > ' ' && raw[end] != ':' && raw[end] != '"' && raw[end] != 0x7f {
			end++
		}
		if end == 0 || end+1 >= len(raw) || raw[end] != ':' || raw[end+1] != '"' {
			return keys, false
		}
		keys = append(keys, raw[:end])
		raw = raw[end+1:]
		closing := 1
		for closing < len(raw) && raw[closing] != '"' {
			if raw[closing] == '\\' {
				closing++
			}
			closing++
		}
		if closing >= len(raw) {
			return keys, false
		}
		raw = raw[closing+1:]
	}
}

func valueToArguments(value string) ([]string, string) {
	var arguments []string
	for value = strings.TrimLeft(value, " "); value != ""; value = strings.TrimLeft(value, " ") {
		end, opener := argumentEnd(value)
		if end == -1 {
			if opener == '{' {
				return nil, unreadable(value, `"union={<key>: <schema>, ...}"`)
			}
			return nil, unreadable(value, "description='<text>'")
		}
		arguments = append(arguments, value[:end])
		value = value[end:]
	}
	return arguments, ""
}

// argumentEnd is -1 when a quote or a brace opened after "=" is not closed.
func argumentEnd(value string) (int, byte) {
	for index := 0; index < len(value); index++ {
		switch value[index] {
		case ' ':
			return index, 0
		case '=':
			opener := charAt(value, index+1)
			closer, found := map[byte]byte{'\'': '\'', '{': '}'}[opener]
			if !found {
				continue
			}
			closing := strings.IndexByte(value[index+2:], closer)
			if closing == -1 {
				return -1, opener
			}
			return index + 2 + closing + 1, opener
		}
	}
	return len(value), 0
}

func valueToFieldOptions(value string) (fieldOptions, string) {
	arguments, err := valueToArguments(value)
	if err != "" {
		return fieldOptions{}, err
	}
	if len(arguments) == 0 || arguments[0] != "field" {
		return fieldOptions{}, `the bidirekt tag must start with "field": write bidirekt:"field" followed by its options`
	}
	values := map[string]string{}
	for _, argument := range arguments[1:] {
		key, optionValue, hasValue := strings.Cut(argument, "=")
		if !hasValue || !slices.Contains(fieldKeys, key) {
			return fieldOptions{}, unreadable(argument, `"<key>=<value>" with key name, type, items, values, optional, nullable, description or union`)
		}
		if _, repeated := values[key]; repeated {
			return fieldOptions{}, fmt.Sprintf(`"%s=" is declared more than once in the bidirekt tag: keep one`, key)
		}
		optionValue, err := optionToValue(key, optionValue, argument)
		if err != "" {
			return fieldOptions{}, err
		}
		values[key] = optionValue
	}
	return valuesToFieldOptions(values)
}

func optionToValue(key, value, argument string) (string, string) {
	quoted := strings.HasPrefix(value, "'") || strings.HasPrefix(value, `"`)
	braced := strings.HasPrefix(value, "{")
	switch {
	case key == "description" && strings.HasPrefix(value, "'"):
		value = strings.TrimSuffix(strings.TrimPrefix(value, "'"), "'")
		return value, descriptionToError(value)
	case key == "description" && !quoted && !braced:
		return value, descriptionToError(value)
	case key == "description":
		return "", unreadable(argument, "description='<text>'")
	case key == "union" && braced:
		if valueToUnion(value) == nil {
			return "", unreadable(argument, `"union={<key>: <schema>, ...}"`)
		}
		return value, ""
	case quoted:
		return "", fmt.Sprintf(`cannot read %q in the bidirekt tag: remove the quotes, only description='<text>' takes them`, argument)
	case braced:
		return "", fmt.Sprintf(`cannot read %q in the bidirekt tag: remove the braces, only union={<key>: <schema>, ...} takes them`, argument)
	}
	switch key {
	case "union":
		return "", unreadable(argument, `"union={<key>: <schema>, ...}"`)
	case "name":
		if value == "" {
			return "", unreadable(argument, `"name=<name>"`)
		}
	case "optional", "nullable":
		if value != "true" && value != "false" {
			return "", unreadable(argument, fmt.Sprintf(`"%s=true" or "%s=false"`, key, key))
		}
	case "type":
		if value != "array" && value != "map" && !schemaNamePattern.MatchString(value) {
			return "", unreadable(argument, typeToFix(value))
		}
	case "items", "values":
		if !schemaNamePattern.MatchString(value) || value == "array" || value == "map" {
			return "", unreadable(argument, fmt.Sprintf(`"%s=" followed by string, integer, float, number, boolean or a schema name`, key))
		}
	}
	return value, ""
}

func typeToFix(value string) string {
	if items, found := strings.CutPrefix(value, "[]"); found {
		return "type=array items=" + items
	}
	if match := mapShortForm.FindStringSubmatch(value); match != nil {
		return "type=map values=" + match[1]
	}
	return `"type=" followed by string, integer, float, number, boolean, a schema name, array or map`
}

func valuesToFieldOptions(values map[string]string) (fieldOptions, string) {
	options := fieldOptions{name: values["name"], description: values["description"]}
	for key, target := range map[string]**bool{"optional": &options.optional, "nullable": &options.nullable} {
		if value, found := values[key]; found {
			flag := value == "true"
			*target = &flag
		}
	}
	typeName, hasType := values["type"]
	_, hasItems := values["items"]
	_, hasValues := values["values"]
	union, hasUnion := values["union"]
	switch {
	case hasType && hasUnion:
		return fieldOptions{}, `declare either "type=" or "union=" in the bidirekt tag, not both`
	case hasItems && typeName != "array":
		return fieldOptions{}, `"items=" is declared without type=array in the bidirekt tag: write type=array or remove "items="`
	case hasValues && typeName != "map":
		return fieldOptions{}, `"values=" is declared without type=map in the bidirekt tag: write type=map or remove "values="`
	case typeName == "array" && !hasItems:
		return fieldOptions{}, `type=array needs "items=" in the bidirekt tag: write the type of its items, as in items=Pet`
	case typeName == "map" && !hasValues:
		return fieldOptions{}, `type=map needs "values=" in the bidirekt tag: write the type of its values, as in values=Pet`
	}
	if hasUnion {
		options.union = valueToUnion(union)
		var schemas []string
		for _, member := range options.union.Members {
			schema := member.Node.Shape.(model.Ref).Schema
			if slices.Contains(schemas, schema) {
				return fieldOptions{}, fmt.Sprintf("the union= lists %q more than once: keep one key for it", schema)
			}
			schemas = append(schemas, schema)
		}
	}
	if !hasType {
		return options, ""
	}
	declared := leafToNode(typeName)
	switch typeName {
	case "array":
		declared = model.Node{Shape: model.Array{Items: leafToNode(values["items"])}}
	case "map":
		declared = model.Node{Shape: model.Map{Values: leafToNode(values["values"])}}
	}
	options.declared = &declared
	return options, ""
}

func leafToNode(name string) model.Node {
	if slices.Contains(primitiveTypes, name) {
		return model.Node{Shape: model.Primitive{Type: name}}
	}
	return model.Node{Shape: model.Ref{Schema: name}}
}

func valueToUnion(value string) *model.Union {
	inner, found := strings.CutPrefix(value, "{")
	if !found {
		return nil
	}
	inner, found = strings.CutSuffix(inner, "}")
	if !found || strings.ContainsAny(inner, "{}") {
		return nil
	}
	union := &model.Union{}
	keys := map[string]bool{}
	for _, member := range strings.Split(inner, ",") {
		match := unionMemberPattern.FindStringSubmatch(strings.TrimSpace(member))
		if match == nil || !schemaNamePattern.MatchString(match[2]) || keys[match[1]] {
			return nil
		}
		keys[match[1]] = true
		union.Members = append(union.Members, model.UnionMember{Key: match[1], Node: model.Node{Shape: model.Ref{Schema: match[2]}}})
	}
	return union
}

func descriptionToError(description string) string {
	if description == "" {
		return `"description" is empty: write a description or leave "description" out`
	}
	if descriptionPattern.MatchString(description) {
		return ""
	}
	quoted, _ := json.Marshal(description)
	return fmt.Sprintf("description %s must be ASCII letters and digits, with single spaces between words: remove or replace the other characters", quoted)
}

func unreadable(argument, usage string) string {
	return fmt.Sprintf("cannot read %q in the bidirekt tag: write %s", argument, usage)
}

func charAt(value string, index int) byte {
	if index >= len(value) {
		return 0
	}
	return value[index]
}
