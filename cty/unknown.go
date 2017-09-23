package cty

func UnknownVal(ty Type) Value {
	return Value{
		v:  nil,
		ty: ty,
	}
}
