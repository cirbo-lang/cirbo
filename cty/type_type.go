package cty

// typeTypeImpl is a typeImpl that allows types to themselves be values within
// the type system.
type typeTypeImpl struct {
	isType
}

// TypeType is a Type whose values are themselves types.
//
// This type is used to allow Cirbo programs to use types in expressions,
// primarily for the purpose of declaring the types of attributes, function
// arguments, etc.
var TypeType Type

func TypeTypeVal(ty Type) Value {
	return Value{
		ty: TypeType,
		v:  ty,
	}
}

func (i typeTypeImpl) Name() string {
	return "Type"
}

func (i typeTypeImpl) Equal(a, b Value) Value {
	at := a.v.(Type)
	bt := b.v.(Type)
	return BoolVal(at.Same(bt))
}

func init() {
	TypeType = Type{typeTypeImpl{}}
}
