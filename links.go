package tjsonapi

import "errors"

var (
	// ErrLinkNotFound is an error object returned when the user tries to
	// access a link that does not exists in a links object.
	ErrLinkNotFound = errors.New("Link not found")

	// ErrLinkUnsupportedType is an error object returned when the user
	// tries to marshal a link of unsupported type (e.g. neither string or
	// Link).
	ErrLinkUnsupportedType = errors.New("Link is of unsupported type")

	// ErrLinkNoMeta is an error object returned when the user tries to
	// retrieve the meta object of a link that doesn't have any. When the link
	// type is Link, it is ensured to have a meta object, although it may be
	// empty.
	ErrLinkNoMeta = errors.New("Link has no meta object")
)

// Link is a struct that represents a link object from the
// <a href="http://jsonapi.org/format/#document-links">JSON API</a>.
type Link struct {
	HRef string                 `json:"href"`
	Meta map[string]interface{} `json:"meta"`
}

// NewLink allocates, initializes and returns a new Link object.
func NewLink() *Link {
	return &Link{
		Meta: make(map[string]interface{}),
	}
}

// Links is a map that associates string values with either *Link or string
// values, effectively representing a links object as defined in the
// <a href="http://jsonapi.org/format/#document-links">JSON API</a>.
type Links map[string]interface{}

// NewLinks allocates a new map and returns it as a Links value.
// The function is the equivalent to calling make(map[string]interface{}).
func NewLinks() Links {
	return make(map[string]interface{})
}

// AddLink adds and associates a string to a given key.
// This function is the equivalent of assigning the string to l[key].
func (l Links) AddLink(key string, link string) {
	l[key] = link
}

// AddLinkObject adds and associates a Link object to a given key.
// This function is the equivalent of assigning the link to l[key].
func (l Links) AddLinkObject(key string, linkObject *Link) {
	l[key] = linkObject
}

// GetLink tries to return the string representation of the link
// associated with the given key. If the link is a Link object, its HRef
// attribute will be returned. If the link is not found or of inappropriate
// type, an empty string will be returned along with the associated error.
func (l Links) GetLink(key string) (string, error) {
	if link, hasKey := l[key]; hasKey {
		switch link.(type) {
		case *Link:
			return link.(*Link).HRef, nil
		case string:
			return link.(string), nil
		default:
			return "", ErrLinkUnsupportedType
		}
	}
	return "", ErrLinkNotFound
}

// GetLinkMeta tries to return the meta object of the link object associated
// with the given key. If the link is either not found, of inappropriate type
// or doesn't have a meta object, a nil map will be returned along with the
// associated error.
func (l Links) GetLinkMeta(key string) (map[string]interface{}, error) {
	if link, hasKey := l[key]; hasKey {
		switch link.(type) {
		case *Link:
			return link.(*Link).Meta, nil
		case string:
			return nil, ErrLinkNoMeta
		default:
			return nil, ErrLinkUnsupportedType
		}
	}
	return nil, ErrLinkNotFound
}
