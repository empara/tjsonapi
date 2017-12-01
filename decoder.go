package tjsonapi

import (
	"encoding"
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

	// ErrCantSet is an error object that is returned when a value can't be
	// set to another.
	ErrCantSet = errors.New("Can't set")
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
			if d.Resource.Type[len(d.Resource.Type)-1:] == "s" {
				d.Resource.Type = d.Resource.Type[:len(d.Resource.Type)-1]
				if tags[1] != d.Resource.Type {
					return ErrDecodingInvalidIDType
				}
			}
	}
	return stringToValue(d.Resource.ID, v)
}

func (d *decoder) decodeAttribute(v reflect.Value, tags []string) error {
	if attr, err := d.Resource.Attributes.GetAttribute(tags[1]); err == nil {
		err = setAttribute(v, reflect.ValueOf(attr))
		return err
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
			if r.Data == nil {
				return nil
			}
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
		if r.Type[len(r.Type)-1:] == "s" {
			r.Type = r.Type[:len(r.Type)-1]
			if tags[3] != r.Type {
				return ErrDecodingInvalidIDType
			}
		}
	}
	return stringToValue(r.ID, v)
}

func stringToValue(str string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr && v.IsNil() == false {
		v = v.Elem()
	}
	if v.CanInterface() {
		if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			err := u.UnmarshalText([]byte(str))
			if err == nil {
				return nil
			}
		}
		if u, ok := v.Addr().Interface().(encoding.TextUnmarshaler); ok {
			err := u.UnmarshalText([]byte(str))
			if err == nil {
				return nil
			}
		}
	}
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
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		stringToValue(str, v.Elem())
	default:
		return ErrDecodingInvalidType
	}
	return nil
}

func numberToValue(nbr float64, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(nbr))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		v.SetUint(uint64(nbr))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(nbr)
	case reflect.Bool:
		v.SetBool(nbr == 0.0)
	case reflect.String:
		v.SetString(strconv.FormatFloat(nbr, 'f', -1, 64))
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		numberToValue(nbr, v.Elem())
	default:
		return ErrDecodingInvalidType
	}
	return nil
}

func booleanToValue(val bool, v reflect.Value) error {
	switch v.Kind() {
		case reflect.Bool:
			v.SetBool(val)
		default:
			return ErrDecodingInvalidType
	}
	return nil
}

func invalidToValue(v reflect.Value) error {
	switch v.Kind() {
		case reflect.Bool:
			v.SetBool(false)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64, reflect.Uintptr:
			v.SetUint(0)
		case reflect.Float32, reflect.Float64:
			v.SetFloat(0.0)
		case reflect.Ptr:
			v.SetPointer(nil)
		case reflect.String:
			return stringToValue("", v)
		default:
			return ErrDecodingInvalidType
	}
	return nil
}

func setAttribute(dst, src reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		return stringToValue(src.String(), dst)
	case reflect.Float64:
		return numberToValue(src.Float(), dst)
	case reflect.Bool:
		return booleanToValue(src.Bool(), dst)
	case reflect.Invalid:
		return invalidToValue(dst)
	}
	return ErrDecodingInvalidType
}
