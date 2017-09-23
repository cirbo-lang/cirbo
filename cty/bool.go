package cty

type boolImpl struct {
	isType
}

var Bool Type

// True is the truthy value of type Bool
var True Value

// True is the falsey value of type Bool
var False Value

func BoolVal(v bool) Value {
	return Value{
		v:  v,
		ty: Bool,
	}
}

func (i boolImpl) Name() string {
	return "Bool"
}

func (i boolImpl) Equal(a, b Value) Value {
	av := a.v.(bool)
	bv := b.v.(bool)
	return BoolVal(av == bv)
}

func (i boolImpl) Not(v Value) Value {
	if v.IsUnknown() {
		return v
	}
	vv := v.v.(bool)
	return BoolVal(!vv)
}

func (i boolImpl) And(a Value, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		return UnknownVal(Bool)
	}
	av := a.v.(bool)
	bv := b.v.(bool)
	return BoolVal(av && bv)
}

func (i boolImpl) Or(a Value, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		return UnknownVal(Bool)
	}
	av := a.v.(bool)
	bv := b.v.(bool)
	return BoolVal(av || bv)
}

func init() {
	Bool = Type{boolImpl{}}
	True = BoolVal(true)
	False = BoolVal(false)
}
