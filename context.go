package tjsonapi

import "errors"

var (
	// ErrContextNotFound is an error object returned when a value is not
	// found in a particular context.
	ErrContextNotFound = errors.New("Value not found in context")
)

// Context is a struct allowing the user to add links and relationships models
// to use with the `jsonapi:"...,context"` tag.
type Context struct {
	Relationships Relationships
	Links         map[string]*Link
}

// NewContext allocates and initializes a new Context object and returns it.
func NewContext() *Context {
	return &Context{
		Relationships: NewRelationships(),
		Links:         make(map[string]*Link),
	}
}
