package tjsonapi

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	// ErrInvalidJSONValue is an error object returned when the value is of
	// a type that can't be marshalled into a JSON value.
	ErrInvalidJSONValue = errors.New("Value is not JSON compatible")
)

var jsonMarshalerType = reflect.TypeOf(new(json.Marshaler)).Elem()

var validJSONKinds = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uintptr,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
	reflect.Interface,
	reflect.Struct,
	reflect.Interface,
	reflect.Struct,
	reflect.Map,
	reflect.Slice,
	reflect.Array,
	reflect.Ptr,
}

// isValidJSONValue checks whether or not the value can be marshalled to JSON.
// Please note that it doesn't check the value recursively, it is up to the
// callee to check for valid values that contains invalid values (e.g.
// structs containing complexes)
func isValidJSONValue(value reflect.Value) error {
	if !value.IsValid() {
		return ErrInvalidJSONValue
	}

	// Check if the value is a supported primitive type
	kind := value.Kind()
	for _, validKinds := range validJSONKinds {
		if kind == validKinds {
			return nil
		}
	}

	// Check if the type or its pointer implements json.Marshaler
	t := value.Type()
	if t.Implements(jsonMarshalerType) ||
		reflect.PtrTo(t).Implements(jsonMarshalerType) {
		return nil
	}
	return ErrInvalidJSONValue
}
