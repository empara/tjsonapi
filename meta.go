package tjsonapi

import (
	"errors"
	"reflect"
)

var (
	// ErrMetaNotFound is an error object returned when the user tries to
	// access a meta member that does not exists in a meta object.
	ErrMetaNotFound = errors.New("Meta not found")
)

// Meta is a map that associates string values with any JSON-compatible type,
// representing a Meta object as defined in the
// <a href="http://jsonapi.org/format/#document-meta">JSON API</a>.
type Meta map[string]interface{}

// NewMeta allocates a new map and returns it as a Meta value.
// This function is the equivalent of doing make(map[string]interface{}).
func NewMeta() Meta {
	return make(map[string]interface{})
}

// AddMeta adds and associates a value to a given key. If the value is not
// a valid JSON value, an error is returned.
// This method is the equivalent of assigning the value to m[key], with
// added sanity checks.
func (m Meta) AddMeta(key string, value interface{}) error {
	err := isValidJSONValue(reflect.ValueOf(value))
	if err != nil {
		return err
	}
	m[key] = value
	return nil
}

// GetMeta tries to return the meta member associated with the given key.
// If the member is not found, this method returns a nil value with an error.
func (m Meta) GetMeta(key string) (interface{}, error) {
	if meta, hasKey := m[key]; hasKey {
		return meta, nil
	}
	return nil, ErrMetaNotFound
}
