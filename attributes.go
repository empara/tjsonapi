package tjsonapi

import (
	"errors"
	"reflect"
)

var (
	// ErrAttributeInvalidKey is an error object returned when a relationships
	// or a links member is added to the attributes object.
	ErrAttributeInvalidKey = errors.New("The attributes object must not " +
		"contain a relationships or links member")

	// ErrAttributeNotFound is an error object returned when the user tries
	// to access an attribute that does not exists in a attributes object.
	ErrAttributeNotFound = errors.New("Attribute not found")
)

// Attributes is a map that associates string values with any JSON-compatible
// type, effectively representing a Attributes object as defined in the
// <a href="http://jsonapi.org/format/#document-resource-object-attributes">
// JSON API</a>.
type Attributes map[string]interface{}

// NewAttributes allocates a new map and returns it as an Attribute value.
// This function is the equivalent of calling make(map[string]interface{}).
func NewAttributes() Attributes {
	return make(map[string]interface{})
}

// AddAttribute adds and associates a value to a given key. It fails and
// returns an error if the key is reserved by the JSON API or if the type of
// the value is not marshalable by Go's JSON package.
// This method is the equivalent of assigning the value to a[key], with
// added sanity checks.
func (a Attributes) AddAttribute(key string, value interface{}) error {
	err := isValidJSONValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	if key == "relationships" || key == "links" {
		return ErrAttributeInvalidKey
	}
	a[key] = value
	return nil
}

// GetAttribute tries to return the attribute associated with the given key.
// If the attribute is not found, this method returns a nil value with an error.
func (a Attributes) GetAttribute(key string) (interface{}, error) {
	if attribute, hasKey := a[key]; hasKey {
		return attribute, nil
	}
	return nil, ErrAttributeNotFound
}
