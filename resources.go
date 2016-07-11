package tjsonapi

import (
	"encoding/json"
	"errors"
)

var (
	// ErrResourcesBadType is an error object returned when the user tries to
	// use a "many" resources object as a "one", and vice-versa.
	ErrResourcesBadType = errors.New("Mismatched types on resources")
)

const (
	// ResourcesOne represents an identifier for resources objects that only
	// contain one resource.
	ResourcesOne = 0

	// ResourcesMany represents an identifier for resources object that contain
	// multiple resource objects.
	ResourcesMany = 1
)

// Resource is a struct that represents a resource object from the
// <a href="http://jsonapi.org/format/#document-resource-objects">JSON API</a>.
type Resource struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
	Links         Links         `json:"links,omitempty"`
	Meta          Meta          `json:"meta,omitempty"`
}

// NewResource allocates and initializes a new Resource object, and returns it.
func NewResource() *Resource {
	return &Resource{
		Attributes:    NewAttributes(),
		Relationships: NewRelationships(),
		Links:         NewLinks(),
		Meta:          NewMeta(),
	}
}

// Resources is a struct representing a resources object from the
// <a href="http://jsonapi.org/format/#document-resource-objects">JSON API</a>.
// To support the ability of having either a single resource or multiple
// resources, the struct contains a type and a resource slice.
type Resources struct {
	Type uint8
	Data []*Resource
}

// NewResourcesOne allocates and initializes a Resources object with
// ResourcesOne as type.
func NewResourcesOne() *Resources {
	return &Resources{
		Type: ResourcesOne,
		Data: make([]*Resource, 1),
	}
}

// NewResourcesMany allocates and initializes a Resources object with
// ResourcesMany as type.
func NewResourcesMany() *Resources {
	return &Resources{
		Type: ResourcesMany,
		Data: make([]*Resource, 0),
	}
}

// SetResource sets the Resource object of a Resources object. If the
// Resources object supports multiple resources instead, an error is returned.
func (r *Resources) SetResource(resource *Resource) error {
	if r.Type == ResourcesOne {
		r.Data[0] = resource
		return nil
	}
	return ErrResourcesBadType
}

// AddResource adds a Resource object to a Resources object. If the Resources
// object doesn't support multiple resources, an error is returned.
func (r *Resources) AddResource(resource *Resource) error {
	if r.Type == ResourcesMany {
		r.Data = append(r.Data, resource)
		return nil
	}
	return ErrResourcesBadType
}

// GetResource tries to return the single Resource object of a Resources
// object. If the Resources object is of "many" type, returns an error.
func (r *Resources) GetResource() (*Resource, error) {
	if r.Type == ResourcesOne {
		return r.Data[0], nil
	}
	return nil, ErrResourcesBadType
}

// GetResources tries to return the internal Resource slice of a Resources
// object. If the Resources object is meant to be used as a single-value
// resource, returns an error instead, as GetResource should be used instead.
func (r *Resources) GetResources() ([]*Resource, error) {
	if r.Type == ResourcesMany {
		return r.Data, nil
	}
	return nil, ErrResourcesBadType
}

// MarshalJSON marshals a Resources object to JSON. This method is needed
// because a Resources object is in fact a multi-type object.<br />
// If there is many resources in the object, the entire slice is marshaled,
// but if there is only one, only the first element is marshaled.
func (r *Resources) MarshalJSON() ([]byte, error) {
	if r.Type == ResourcesOne {
		return json.Marshal(r.Data[0])
	}
	return json.Marshal(r.Data)
}
