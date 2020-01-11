package client

// Filter represents a Skylark request filter
// It is used to constrain a request by a field's value
type Filter struct {
	c     constraint
	value string
}

type constraint string

const (
	// GreaterThan contrains to fields that are greater than a given value
	GreaterThan = constraint("gt")
	// LessThan contrains to fields that are less than a given value
	LessThan = constraint("lt")
	// Equals contrains to fields that equal a given value
	Equals = constraint("")
)

// NewFilter creates a new filter with a given constraint and value
func NewFilter(c constraint, value string) *Filter {
	return &Filter{c: c, value: value}
}
