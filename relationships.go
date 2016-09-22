package tjsonapi

import (
	"encoding/json"
	"errors"
)

var (
	// ErrResourceLinkageBadType is an error object returned when the user
	// tries to do a to-many operation on a to-one resource linkage and
	// vice-versa.
	ErrResourceLinkageBadType = errors.New("Mismatched types on resource " +
		"linkage")
)

const (
	// ResourceLinkageToOne represents an identifier for to-one resource
	// linkages to use as type.
	ResourceLinkageToOne = 0

	// ResourceLinkageToMany represents an identifier for to-many resource
	// linkages to use as type.
	ResourceLinkageToMany = 1
)

// ResourceIdentifier is a struct that represents a resource linkage object
// from the
// <a href="http://jsonapi.org/format/#document-resource-object-linkage">JSON
// API</a>.
type ResourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Meta Meta   `json:"meta,omitempty"`
}

// NewResourceIdentifier allocates and initializes a new ResourceIdentifier
// object.
func NewResourceIdentifier() *ResourceIdentifier {
	return &ResourceIdentifier{
		Meta: NewMeta(),
	}
}

// ResourceLinkage is a struct representing a ResourceLinkage object from the
// <a href="http://jsonapi.org/format/#document-resource-object-linkage">JSON
// API</a>. To support both to-one and to-many linkages, the struct contains
// a type and a resource identifier slice.
type ResourceLinkage struct {
	Type uint8
	Data []*ResourceIdentifier
}

// NewResourceLinkageToOne allocates and initializes a ResourceLinkage object
// with ResourceLinkageToOne as type.
func NewResourceLinkageToOne() *ResourceLinkage {
	return &ResourceLinkage{
		Type: ResourceLinkageToOne,
		Data: make([]*ResourceIdentifier, 1),
	}
}

// NewResourceLinkageToMany allocates and initializes a ResourceLinkage object
// with ResourceLinkageToMany as type.
func NewResourceLinkageToMany() *ResourceLinkage {
	return &ResourceLinkage{
		Type: ResourceLinkageToMany,
		Data: make([]*ResourceIdentifier, 0),
	}
}

// AddResourceIdentifier adds a ResourceIdentifier object to a to-many
// ResourceLinkage. If the ResourceLinkage is to-one, does nothing and returns
// an error.
func (l *ResourceLinkage) AddResourceIdentifier(r *ResourceIdentifier) error {
	if l.Type == ResourceLinkageToMany {
		l.Data = append(l.Data, r)
		return nil
	}
	return ErrResourceLinkageBadType
}

// SetResourceIdentifier sets the ResourceIdentifier object of a to-one
// ResourceLinkage. If the ResourceLinkage is to-many, does nothing and returns
// an error.
func (l *ResourceLinkage) SetResourceIdentifier(r *ResourceIdentifier) error {
	if l.Type == ResourceLinkageToOne {
		l.Data[0] = r
		return nil
	}
	return ErrResourceLinkageBadType
}

// GetResourceIdentifiers tries to return the internal ResourceIdentifier slice
// of a ResourceLinkage object. Although it's possible to get a single-item
// slice for to-one ResourceLinkage, this function returns an error when the
// user attempts it, because it is considered as a to-many operation.
func (l *ResourceLinkage) GetResourceIdentifiers() ([]*ResourceIdentifier,
	error) {
	if l.Type == ResourceLinkageToMany {
		return l.Data, nil
	}
	return nil, ErrResourceLinkageBadType
}

// GetResourceIdentifier tries to return the ResourceIdentifier object contained
// in a ResourceLinkage object. If this function is used in a to-many context,
// it does nothing and returns an error.
func (l *ResourceLinkage) GetResourceIdentifier() (*ResourceIdentifier, error) {
	if l.Type == ResourceLinkageToOne {
		return l.Data[0], nil
	}
	return nil, ErrResourceLinkageBadType
}

// MarshalJSON marshals a ResourceLinkage to JSON. This method is needed because
// a ResourceLinkage object is in fact a multi-type object.<br />
// If the ResourceLinkage is to-one, the value of the first and unique element
// of the underlying slice is marshaled. If the ResourceLinkage is to-many,
// the entire slice is then marshaled.
func (l *ResourceLinkage) MarshalJSON() ([]byte, error) {
	if l.Type == ResourceLinkageToOne {
		return json.Marshal(l.Data[0])
	}
	return json.Marshal(l.Data)
}

// UnmarshalJSON unmarshals JSON data to a ResourceLinkage object.
// Can fail if the JSON data does not represent an object or an array.
func (l *ResourceLinkage) UnmarshalJSON(data []byte) error {
	var linkage ResourceIdentifier
	if err := json.Unmarshal(data, &linkage); err == nil {
		l.Type = ResourceLinkageToOne
		l.Data = []*ResourceIdentifier{&linkage}
		return nil
	}

	if err := json.Unmarshal(data, &l.Data); err == nil {
		l.Type = ResourceLinkageToMany
		return nil
	}

	return ErrResourceLinkageBadType
}

// Relationship is a struct that represents a relationship object from the
// <a href="http://jsonapi.org/format/#document-resource-object-relationships">
// JSON API</a>.
type Relationship struct {
	Links Links            `json:"links,omitempty"`
	Data  *ResourceLinkage `json:"data,omitempty"`
	Meta  Meta             `json:"meta,omitempty"`
}

// NewRelationship allocates, initializes and returns a Relationship object.
func NewRelationship() *Relationship {
	return &Relationship{
		Links: NewLinks(),
		Meta:  NewMeta(),
	}
}

// Relationships is a map that associates string values with *Relationship
// values, effectively representing a relationships object as defined in the
// <a href="http://jsonapi.org/format/#document-resource-object-relationships">
// JSON API</a>.
type Relationships map[string]*Relationship

// NewRelationships allocates a new map and returns it as a Relationships value.
// This function is the equivalent to calling make(map[string]*Relationship).
func NewRelationships() Relationships {
	return make(map[string]*Relationship)
}
