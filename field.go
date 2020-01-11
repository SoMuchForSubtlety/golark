package client

import (
	"fmt"
	"net/url"
)

// Field represents a Skylark request field
type Field struct {
	Name       string
	IsIncluded bool
	IsExpanded bool
	SubFields  map[string]*Field
	filters    []*Filter
}

// NewField creates a new field
func NewField(name string) *Field {
	return &Field{Name: name, SubFields: make(map[string]*Field), IsIncluded: true}
}

func (f *Field) apply(v url.Values) url.Values {
	if f.IsIncluded {
		v = addValue(v, "fields", f.Name)
	}
	if f.IsExpanded {
		v = addValue(v, "fields_to_expand", f.Name)
	}
	for _, filter := range f.filters {
		var key string
		if filter.c == "" {
			key = f.Name
		} else {
			key = fmt.Sprintf("%s__%s", f.Name, filter.c)
		}
		v.Add(key, filter.value)
	}
	for _, field := range f.SubFields {
		v = field.apply(v)
	}
	return v
}

// WithSubField expands a field and adds the given field to the list of filds to be returned.
// Only use this if the field is a reference to a different object!
func (f *Field) WithSubField(subField *Field) *Field {
	f.IsExpanded = true
	subField.adjustName(f.Name)
	f.SubFields[subField.Name] = subField
	return f
}

func (f *Field) adjustName(parentName string) {
	f.Name = fmt.Sprintf("%s__%s", parentName, f.Name)
	for _, field := range f.SubFields {
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
	subField.IsExpanded = true
	subField.IsIncluded = false
	subField.adjustName(f.Name)
	f.SubFields[subField.Name] = subField
	return f
}
