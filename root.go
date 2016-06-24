package tjsonapi

// Root is a struct that represents the top-level object of a
// <a href="http://jsonapi.org/format/#document-top-level">JSON API</a>
// document.
type Root struct {
	Data *Resources `json:"data,omitempty"`
	Meta Meta       `json:"meta,omitempty"`
}

// NewRoot allocates a new Root object. Equivalent to new(Root).
func NewRoot() *Root {
	return new(Root)
}
