package golark

import (
	"fmt"
	"net/url"
)

// Field represents a Skylark request field
type Field struct {
	name       string
	isIncluded bool
	isExpanded bool
	subFields  map[string]*Field
	filters    []*Filter
}

// NewField creates a new field
func NewField(name string) *Field {
	return &Field{name: name, subFields: make(map[string]*Field), isIncluded: true}
}

func (f *Field) apply(v url.Values) url.Values {
	if f.isIncluded {
		v = addValue(v, "fields", f.name)
	}
	if f.isExpanded {
		v = addValue(v, "fields_to_expand", f.name)
	}
	for _, filter := range f.filters {
		var key string
		if filter.c == "" {
			key = f.name
		} else {
			key = fmt.Sprintf("%s__%s", f.name, filter.c)
		}
		v.Add(key, filter.value)
	}
	for _, field := range f.subFields {
		v = field.apply(v)
	}
	return v
}

// WithSubField expands a field and adds the given field to the list of filds to be returned.
// Only use this if the field is a reference to a different object!
func (f *Field) WithSubField(subField *Field) *Field {
	f.isExpanded = true
	subField.adjustName(f.name)
	f.subFields[subField.name] = subField
	return f
}

func (f *Field) adjustName(parentName string) {
	f.name = fmt.Sprintf("%s__%s", parentName, f.name)
	for _, field := range f.subFields {
		field.adjustName(parentName)
	}
}

// WithFilter applies a fielter to the field.
func (f *Field) WithFilter(filter *Filter) *Field {
	f.filters = append(f.filters, filter)
	return f
}

// Expand expands a field without explicitly listing it as a field to return.
// This is usefult if you want to return all fields.
func (f *Field) Expand(subField *Field) *Field {
	subField.isExpanded = true
	subField.isIncluded = false
	subField.adjustName(f.name)
	f.subFields[subField.name] = subField
	return f
}
