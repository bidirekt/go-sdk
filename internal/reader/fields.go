package reader

import (
	"cmp"
	"go/types"
	"slices"
	"strings"
	"unicode"
)

type directField struct {
	variable  *types.Var
	index     int
	tag       structTag
	jsonName  string
	omitEmpty bool
	omitZero  bool
	quoted    bool
	skipped   bool
	embedded  *types.Struct
	embedName *types.TypeName
}

type jsonField struct {
	name   string
	tagged bool
	index  []int
	field  *directField
}

type structFields struct {
	visible []jsonField
	marked  bool
}

type embeddedLevel struct {
	fields *types.Struct
	index  []int
}

func (this *reading) structToDirectFields(fields *types.Struct) []*directField {
	if direct, found := this.directFields[fields]; found {
		return direct
	}
	direct := make([]*directField, 0, fields.NumFields())
	for index := range fields.NumFields() {
		direct = append(direct, this.varToDirectField(fields.Field(index), index, fields.Tag(index)))
	}
	this.directFields[fields] = direct
	return direct
}

func (this *reading) varToDirectField(variable *types.Var, index int, raw string) *directField {
	field := &directField{variable: variable, index: index, tag: rawToStructTag(raw)}
	for _, message := range field.tag.warnings {
		this.addWarningAt(variable.Pos(), "%s", message)
	}
	for _, message := range field.tag.errors {
		this.addErrorAt(variable.Pos(), "%s", message)
	}
	name, options, _ := strings.Cut(field.tag.jsonTag, ",")
	if isValidJSONName(name) {
		field.jsonName = name
	}
	written := deref(variable.Type())
	_, isStruct := written.Underlying().(*types.Struct)
	switch {
	case field.tag.jsonTag == "-":
		field.skipped = true
		if field.tag.marked {
			this.addErrorAt(variable.Pos(), `field %q has json:"-", so encoding/json does not write it: remove json:"-" or the bidirekt tag`, variable.Name())
		}
		return field
	case !variable.Exported() && (!variable.Embedded() || !isStruct):
		field.skipped = true
		if field.tag.marked {
			this.addErrorAt(variable.Pos(), "field %q is not exported, so encoding/json does not write it: export it or remove the bidirekt tag", variable.Name())
		}
		return field
	}
	for option := range strings.SplitSeq(options, ",") {
		field.omitEmpty = field.omitEmpty || option == "omitempty"
		field.omitZero = field.omitZero || option == "omitzero"
		field.quoted = field.quoted || (option == "string" && isQuotable(written))
	}
	if variable.Embedded() && field.jsonName == "" && isStruct {
		field.embedded = written.Underlying().(*types.Struct)
		if named, isNamed := types.Unalias(written).(*types.Named); isNamed {
			field.embedName = named.Obj()
		}
		if field.tag.marked {
			this.addErrorAt(variable.Pos(), "embedded struct %q is not a field of its own, its marked fields are promoted: remove its bidirekt tag, or give it a json name to make it a field", variable.Name())
		}
	}
	return field
}

// structToFields follows encoding/json: the fields of embedded structs without a json name are promoted,
// and of the fields with the same name the shallowest wins, then the one with a json name; when neither
// wins, encoding/json writes none of them.
func (this *reading) structToFields(fields *types.Struct) structFields {
	if found, cached := this.jsonFields[fields]; cached {
		return found
	}
	var all []jsonField
	marked := false
	visited := map[*types.Struct]bool{}
	count := map[*types.Struct]int{}
	next := []embeddedLevel{{fields: fields}}
	for len(next) > 0 {
		current := next
		next = nil
		nextCount := map[*types.Struct]int{}
		for _, level := range current {
			if visited[level.fields] {
				continue
			}
			visited[level.fields] = true
			for _, field := range this.structToDirectFields(level.fields) {
				marked = marked || field.tag.marked
				if field.skipped {
					continue
				}
				index := append(slices.Clone(level.index), field.index)
				if field.embedded == nil {
					name := field.jsonName
					if name == "" {
						name = field.variable.Name()
					}
					all = append(all, jsonField{name: name, tagged: field.jsonName != "", index: index, field: field})
					if count[level.fields] > 1 {
						all = append(all, all[len(all)-1])
					}
					continue
				}
				if field.embedName != nil {
					this.embeddedNames[field.embedName] = true
				}
				nextCount[field.embedded]++
				if nextCount[field.embedded] == 1 {
					next = append(next, embeddedLevel{fields: field.embedded, index: index})
				}
			}
		}
		count = nextCount
	}
	result := structFields{visible: this.fieldsToDominant(all), marked: marked}
	this.jsonFields[fields] = result
	return result
}

func (this *reading) fieldsToDominant(all []jsonField) []jsonField {
	slices.SortFunc(all, func(a, b jsonField) int {
		return cmp.Or(
			strings.Compare(a.name, b.name),
			cmp.Compare(len(a.index), len(b.index)),
			compareTagged(a, b),
			slices.Compare(a.index, b.index),
		)
	})
	var visible []jsonField
	for start := 0; start < len(all); {
		end := start + 1
		for end < len(all) && all[end].name == all[start].name {
			end++
		}
		group := all[start:end]
		if len(group) > 1 && len(group[0].index) == len(group[1].index) && group[0].tagged == group[1].tagged {
			this.addConflict(group)
		} else {
			visible = append(visible, group[0])
		}
		start = end
	}
	slices.SortFunc(visible, func(a, b jsonField) int { return slices.Compare(a.index, b.index) })
	return visible
}

func (this *reading) addConflict(group []jsonField) {
	for _, each := range group {
		if each.field.tag.marked && len(each.index) == len(group[0].index) {
			this.addErrorAt(each.field.variable.Pos(), "json name %q is written by two fields at the same depth of embedded structs, so encoding/json writes neither: give one of them another json name", each.name)
		}
	}
}

func compareTagged(a, b jsonField) int {
	switch {
	case a.tagged == b.tagged:
		return 0
	case a.tagged:
		return -1
	}
	return 1
}

func isValidJSONName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c) && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isQuotable(written types.Type) bool {
	basic, isBasic := written.Underlying().(*types.Basic)
	return isBasic && basic.Info()&(types.IsBoolean|types.IsNumeric|types.IsString) != 0 && basic.Info()&types.IsComplex == 0
}

func deref(written types.Type) types.Type {
	if pointer, isPointer := types.Unalias(written).(*types.Pointer); isPointer {
		return pointer.Elem()
	}
	return written
}
