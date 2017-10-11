package cbty

import (
	"golang.org/x/text/unicode/norm"
)

type stringImpl struct {
	isType
}

var String Type

func StringVal(s string) Value {
	return Value{
		v:  norm.NFC.String(s),
		ty: String,
	}
}

func (i stringImpl) Name() string {
	return "String"
}

func (i stringImpl) Equal(a, b Value) Value {
	av := a.v.(string)
	bv := b.v.(string)
	return BoolVal(av == bv)
}

func (i stringImpl) CanConcat(o Type) bool {
	_, isString := o.impl.(stringImpl)
	return isString
}

func (i stringImpl) Concat(a Value, b Value) Value {
	if a.IsUnknown() || b.IsUnknown() {
		return UnknownVal(String)
	}

	return StringVal(a.v.(string) + b.v.(string))
}

func init() {
	String = Type{stringImpl{}}
}
