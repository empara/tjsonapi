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
		Context: c,
	}

	root := new(Root)
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Struct:
		root.Data = NewResourcesOne()
		e.Resource = NewResource()
		err := e.marshalStruct(v)
		if err != nil {
			return nil, err
		}
		root.Data.SetResource(e.Resource)
	case reflect.Array, reflect.Slice:
		root.Data = NewResourcesMany()
		for it := 0; it < v.Len(); it++ {
			e.Resource = NewResource()
			err := e.marshalStruct(v.Index(it))
			if err != nil {
				return nil, err
			}
			root.Data.AddResource(e.Resource)
		}
	default:
		return nil, ErrEncodingInvalidType
	}
	return root, nil
}

type encoder struct {
	Context           *Context
	Resource          *Resource
	RelationshipCount int
}

func (e *encoder) marshalStruct(v reflect.Value) error {
	t := v.Type()
	for it := 0; it < t.NumField(); it++ {
		f := t.Field(it)
		tags := strings.Split(f.Tag.Get("jsonapi"), ",")

		var err error
		switch tags[0] {
		case TagIdentifier:
			err = e.encodeIdentifier(v.Field(it), tags)
		case TagAttribute:
			err = e.encodeAttribute(v.Field(it), tags)
		case TagRelationship:
			err = e.encodeRelationship(v.Field(it), tags)
		}
		if err != nil {
			return err
		}
	}
	return nil
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
	return e.Resource.Attributes.AddAttribute(tags[1], v.Interface())
}

func (e *encoder) encodeRelationship(v reflect.Value, tags []string) error {
	if len(tags) < 2 {
		return ErrEncodingInvalidTag
	}

	// Only `jsonapi:"relationship,[key]"` has a 2-value tag.
	if len(tags) == 2 {
		var r Relationship
		if v.Kind() == reflect.Ptr &&
			v.Type() == reflect.PtrTo(reflect.TypeOf(r)) {
			r = v.Elem().Interface().(Relationship)
		} else if v.Type() == reflect.TypeOf(r) {
			r = v.Interface().(Relationship)
		} else {
			return ErrEncodingInvalidType
		}
		e.Resource.Relationships[tags[1]] = &r
	} else {
		if len(tags) < 3 {
			return ErrEncodingInvalidTag
		}
		switch tags[2] {
		case TagRelationshipContext:
			rPtr := e.Context.Relationships[tags[1]]
			if rPtr == nil {
				return ErrContextNotFound
			}
			r := *rPtr
			populateStruct(reflect.ValueOf(r), v)
			e.Resource.Relationships[tags[1]] = &r
		case TagRelationshipLink:
			if v.Kind() != reflect.String {
				return ErrEncodingInvalidType
			}
			r := NewRelationship()
			r.Links.AddLink("self", v.String())
			e.Resource.Relationships[tags[1]] = r
		case TagRelationshipData:
			if len(tags) < 4 {
				return ErrEncodingInvalidTag
			}
			var err error
			r := NewRelationship()
			resource := NewResourceIdentifier()
			resource.ID, err = valueToString(v)
			resource.Type = tags[3]
			if err != nil {
				return ErrEncodingInvalidType
			}
			r.Data = NewResourceLinkageToOne()
			r.Data.SetResourceIdentifier(resource)
			e.Resource.Relationships[tags[1]] = r
		}
	}
	return nil
}

func (e *encoder) encodeLink(v reflect.Value, tags []string) error {
	if len(tags) < 2 {
		return ErrEncodingInvalidTag
	}

	// Only `jsonapi:"link,[key]"` has a 2-value tag.
	if len(tags) == 2 {
		str, err := valueToString(v)
		if err != nil {
			return err
		}
		e.Resource.Links.AddLink(tags[1], str)
	} else {
		if len(tags) < 3 {
			return ErrEncodingInvalidTag
		}
		switch tags[2] {
		case "context":
			lPtr := e.Context.Links[tags[1]]
			if lPtr == nil {
				return ErrContextNotFound
			}
			l := *lPtr
			populateStruct(reflect.ValueOf(l), v)
			e.Resource.Links[tags[1]] = v
		}
	}
	return nil
}

func (e *encoder) encodeMeta(v reflect.Value, tags []string) error {
	if len(tags) < 2 {
		return ErrEncodingInvalidTag
	}
	e.Resource.Meta[tags[1]] = v.Interface()
	return nil
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

// populateStruct will browse a given value v and replace every field marked
// with `jsonapi:"value"` with the value i. Is recursive.
func populateStruct(v, i reflect.Value) error {
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for it := 0; it < t.NumField(); it++ {
			f := t.Field(it)
			err := populateField(f, v.Field(it), i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func populateField(f reflect.StructField, v, i reflect.Value) error {
	// If the type of the value is a pointer or a struct, we can go deeper
	switch v.Kind() {
	case reflect.Ptr:
		return populateField(f, v.Elem(), i)
	case reflect.Struct:
		return populateStruct(v, i)
	}

	if f.Tag.Get("jsonapi") == TagValue {
		if v.CanSet() && i.Type().AssignableTo(v.Type()) {
			v.Set(i)
		}
		return ErrEncodingInvalidType
	}
	return nil
}
