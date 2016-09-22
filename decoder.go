package tjsonapi

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrDecodingInvalidType is an error object that is returned when the
	// interface sent to Unmarshal mismatches the Resources object (e.g. not
	// a pointer when Resources is "One" or not a slice when Resources is
	// "many").
	ErrDecodingInvalidType = errors.New("Interface has invalid type")

	// ErrDecodingInvalidIDType is an error object that is returned when the
	// type specified in the JSONAPI identifier tag doesn't match the one in
	// the decoded data.
	ErrDecodingInvalidIDType = errors.New("Struct has invalid identifier type")

	// ErrDecodingInvalidTag is an error object that is returned when the tag
	// for a field is invalid (e.g. when is a tag is unknown, or when the tag
	// doesn't contain enough sub-tags).
	ErrDecodingInvalidTag = errors.New("Invalid JSONAPI tag")
)

// Unmarshal fills up an interface from a JSONAPI root.
// This function is equivalent to creating a blank Context and unmarshaling
// the root with it.
// See Context.Unmarshal for more details.
func Unmarshal(r *Root, v interface{}) error {
	c := new(Context)
	return c.Unmarshal(r, v)
}

// Unmarshal fills up an interface from a JSONAPI root, using c as the Context.
func (c *Context) Unmarshal(r *Root, i interface{}) error {
	d := &decoder{
		Context: c,
	}

	if r.Data.Type == ResourcesOne {
		v := reflect.ValueOf(i)
		if v.Kind() == reflect.Ptr {
			d.Resource = r.Data.Data[0]
			return d.unmarshalResource(v.Elem())
		}
	} else if r.Data.Type == ResourcesMany {
		v := reflect.ValueOf(i)
		if v.Kind() == reflect.Slice {
			v.SetLen(0)
			vType := v.Type().Elem()
			for _, resource := range r.Data.Data {
				vElem := reflect.New(vType)
				d.Resource = resource
				err := d.unmarshalResource(vElem)
				if err != nil {
					return err
				}
				reflect.Append(v, vElem)
			}
			return nil
		}
	}
	return ErrDecodingInvalidType
}

type decoder struct {
	Context  *Context
	Resource *Resource
}

func (d *decoder) unmarshalResource(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := strings.Split(f.Tag.Get("jsonapi"), ",")

		var err error
		switch tags[0] {
		case TagIdentifier:
			err = d.decodeIdentifier(v.Field(i), tags)
		case TagAttribute:
			err = d.decodeAttribute(v.Field(i), tags)
		case TagRelationship:
			err = d.decodeRelationship(v.Field(i), tags)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeIdentifier(v reflect.Value, tags []string) error {
	if tags[1] != d.Resource.Type {
		return ErrDecodingInvalidIDType
	}
	return stringToValue(d.Resource.ID, v)
}

func (d *decoder) decodeAttribute(v reflect.Value, tags []string) error {
	if attr, err := d.Resource.Attributes.GetAttribute(tags[1]); err == nil {
		v.Set(reflect.ValueOf(attr))
	}
	return nil
}

// This method only supports "data" relationships. For now.
func (d *decoder) decodeRelationship(v reflect.Value, tags []string) error {
	switch tags[2] {
	case TagRelationshipData:
		if len(tags) < 4 {
			return ErrDecodingInvalidTag
		}
		if r, hasKey := d.Resource.Relationships[tags[1]]; hasKey {
			if r.Data.Type == ResourceLinkageToOne {
				return d.decodeResourceIdentifier(v, r.Data.Data[0], tags)
			} else if r.Data.Type == ResourceLinkageToMany {
				for it := 0; it < len(r.Data.Data); it++ {
					vType := v.Type().Elem()
					vElem := reflect.New(vType)
					err := d.decodeResourceIdentifier(vElem,
						r.Data.Data[it], tags)
					if err != nil {
						return err
					}
					v.Set(reflect.Append(v, vElem.Elem()))
				}
			} else {
				return ErrDecodingInvalidType
			}
		}
	}
	return nil
}

func (d *decoder) decodeResourceIdentifier(v reflect.Value,
	r *ResourceIdentifier, tags []string) error {
	if tags[3] != r.Type {
		return ErrDecodingInvalidIDType
	}
	return stringToValue(r.ID, v)
}

func stringToValue(str string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		nb, _ := strconv.ParseInt(str, 10, 64)
		v.SetInt(nb)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		nb, _ := strconv.ParseUint(str, 10, 64)
		v.SetUint(nb)
	case reflect.Float32, reflect.Float64:
		nb, _ := strconv.ParseFloat(str, 64)
		v.SetFloat(nb)
	case reflect.Bool:
		boolean, _ := strconv.ParseBool(str)
		v.SetBool(boolean)
	case reflect.String:
		v.SetString(str)
	case reflect.Ptr:
		stringToValue(str, v.Elem())
	default:
		return ErrEncodingInvalidType
	}
	return nil
}
