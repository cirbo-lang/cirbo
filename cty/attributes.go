package cty

// typeWithAttributes is an interface implemented by typeImpls that can
// support attribute lookup.
type typeWithAttributes interface {
	// GetAttr returns the value of the attribute with the given name, or
	// NilValue if such an attribute is not defined.
	GetAttr(name string) Value
}

// staticAttributes can be embedded into a type to implement typeWithAttributes
// using a fixed map of attribute values.
type staticAttributes map[string]Value

func (a staticAttributes) GetAttr(name string) Value {
	return a[name]
}
