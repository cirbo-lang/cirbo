package cty

func UnknownVal(ty Type) Value {
	return Value{
		v:  nil,
		ty: ty,
	}
}

type unknownTypeImpl [0]*unknownTypeImpl

func (i unknownTypeImpl) typeSigil() isType {
	return isType{}
}

func (i unknownTypeImpl) Name() string {
	return "<unknown type>"
}

func (i unknownTypeImpl) GoString() string {
	return "cty.UnknownType"
}

func (i unknownTypeImpl) Equal(a, b Value) Value {
	return UnknownVal(Bool)
}

// PlaceholderVal is an unknown value whose type is also unknown. This can
// be used as a placeholder where a valid value is required but no specific
// value is appropriate. It implements all operations with itself as the
// result.
var PlaceholderVal Value

// FIXME: PlaceholderVal doesn't actually support all of the operations yet.

func init() {
	PlaceholderVal = UnknownVal(Type{
		impl: unknownTypeImpl{},
	})
}
