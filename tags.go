package tjsonapi

const (
	// TagIdentifier is the top-level tag used to define a value as an
	// identifier.
	TagIdentifier = "identifier"

	// TagAttribute is the top-level tag used to define a value as an
	// attribute.
	TagAttribute = "attribute"

	// TagRelationship is the top-level tag used to mark a value as a
	// relationship.
	TagRelationship = "relationship"

	// TagLink is the top-level tag used to define a value as a link.
	TagLink = "link"

	// TagMeta is the top-level tag used to define a value as a part of the
	// meta object.
	TagMeta = "meta"

	// TagValue is the tag used when populating structs of a context. Member
	// marked with this tag can and will be set to a specific value when
	// using contexted members.
	TagValue = "value"

	// TagRelationshipContext is the sub-tag used to define a value as a
	// context relationship.
	TagRelationshipContext = "context"

	// TagRelationshipLink is the sub-tag used to define a value as a link
	// relationship.
	TagRelationshipLink = "link"

	// TagRelationshipData is the sub-tag used to define a value as a resource
	// linkage relationship.
	TagRelationshipData = "data"

	// TagLinkContext is the sub-tag used to define a value as a context link.
	TagLinkContext = "context"
)
