package tjsonapi

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrEncodingInvalidType is an error object that is returned when the
	// value send to marshal is not a struct.
	ErrEncodingInvalidType = errors.New("Encoding JSONAPI from invalid type")

	// ErrEncodingInvalidTag is an error object that is returned when the
	// tag for a field is invalid (e.g. when a tag is unknown, or when
	// the tag doesn't contain enough sub-tags).
	ErrEncodingInvalidTag = errors.New("Invalid JSONAPI tag")
)

// Marshal returns a JSON-marshalable root for the given interface.
// This function is equivalent to creating a blank Context and marshaling the
// interface with it.
// See Context.Marshal for more details.
func Marshal(v interface{}) (*Root, error) {
	c := new(Context)
	return c.Marshal(v)
}

// Marshal returns a JSON-marshalable root for the given interface, using c
// as the Context.
func (c *Context) Marshal(i interface{}) (*Root, error) {
	e := &encoder{
		Context:  c,
		Resource: NewResource(),
	}

	v := reflect.ValueOf(i)
	t := v.Type()
	if v.Kind() == reflect.Struct {
		for it := 0; it < t.NumField(); it++ {
			f := t.Field(it)
			tags := strings.Split(f.Tag.Get("jsonapi"), ",")

			var err error
			switch tags[0] {
			case TagIdentifier:
				err = e.encodeIdentifier(v.Field(it), tags)
			case TagAttribute:
				err = e.encodeAttribute(v.Field(it), tags)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	root := new(Root)
	root.Data = NewResourcesOne()
	root.Data.SetResource(e.Resource)
	return root, nil
}

type encoder struct {
	Context  *Context
	Resource *Resource
}

func (e *encoder) encodeIdentifier(v reflect.Value, tags []string) (err error) {
	if len(tags) < 2 {
		return ErrEncodingInvalidTag
	}
	e.Resource.ID, err = valueToString(v)
	e.Resource.Type = tags[1]
	return
}

func (e *encoder) encodeAttribute(v reflect.Value, tags []string) error {
	if len(tags) < 2 {
		return ErrEncodingInvalidTag
	}
	return e.Resource.Attributes.AddAttribute(tags[1], v)
}

func valueToString(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "", ErrEncodingInvalidType
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Ptr:
		return valueToString(v.Elem())
	default:
		return "", ErrEncodingInvalidType
	}
}
