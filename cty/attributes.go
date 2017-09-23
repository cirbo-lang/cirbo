package cty

// typeWithAttributes is an interface implemented by typeImpls that can
// support attribute lookup.
type typeWithAttributes interface {
	// GetAttr returns the value of the attribute with the given name, or
	// NilValue if such an attribute is not defined.
	//
	// If the given value is not known, the result is an unknown value of
	// the attribute's type.
	GetAttr(val Value, name string) Value
}

// staticAttributes can be embedded into a type to implement typeWithAttributes
// using a fixed map of attribute values.
type staticAttributes map[string]Value

func (a staticAttributes) GetAttr(val Value, name string) Value {
	v := a[name]
	if v == NilValue {
		return NilValue
	}

	if !val.IsKnown() {
		return UnknownVal(v.Type())
	}

	return v
}
