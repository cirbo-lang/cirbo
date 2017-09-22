package cty

type stringImpl struct {
	isType
}

var String Type

func StringVal(s string) Value {
	return Value{
		v:  s,
		ty: String,
	}
}

func (i stringImpl) Name() string {
	return "String"
}

func (i stringImpl) CanConcat(o Type) bool {
	_, isString := o.impl.(stringImpl)
	return isString
}

func (i stringImpl) Concat(a Value, b Value) Value {
	return StringVal(a.v.(string) + b.v.(string))
}

func init() {
	String = Type{stringImpl{}}
}
