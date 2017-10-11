package cbty

import (
	"bytes"
	"fmt"
	"sort"
)

type objectImpl struct {
	isType
	atys map[string]Type
}

// EmptyObject is an alias for an object type with no attributes at all.
var EmptyObject = Object(map[string]Type{})

// EmptyObjectVal is the only known value of type EmptyObject.
var EmptyObjectVal Value

// Object creates a new object type with the given attribute types.
//
// Object is a generic data structure with attributes. It has no special
// meaning, and more meaningful types that themselves have attributes are
// not considered subtypes of this type.
func Object(atys map[string]Type) Type {
	if atys == nil {
		panic("attempt to create Object type with nil attribute type map")
	}

	return Type{objectImpl{atys: atys}}
}

// ObjectVal creates a value of an object type constructed from the types
// of the given attribute values.
//
// Although this function does not enforce it, the language internals assume
// that attribute names will always be valid identifiers in the language
// syntax. An object with invalid attribute names will cause undefined
// behavior.
func ObjectVal(attrs map[string]Value) Value {
	if len(attrs) == 0 {
		return EmptyObjectVal
	}
	atys := map[string]Type{}
	rawVs := map[string]interface{}{}
	for n, v := range attrs {
		atys[n] = v.ty
		rawVs[n] = v.v
	}
	return Value{
		v:  rawVs,
		ty: Object(atys),
	}
}

func (i objectImpl) Name() string {
	var buf bytes.Buffer
	buf.WriteString("Object(")
	first := true
	names := make([]string, 0, len(i.atys))
	for name := range i.atys {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		ty := i.atys[name]
		if !first {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s=%s", name, ty.Name())
		first = false
	}
	buf.WriteString(")")
	return buf.String()
}

func (i objectImpl) GoString() string {
	if len(i.atys) == 0 {
		return "cty.EmptyObject"
	}

	return fmt.Sprintf("cty.Object(%#v)", i.atys)
}

func (i objectImpl) Same(o Type) bool {
	oi, isObj := o.impl.(objectImpl)
	if !isObj {
		return false
	}

	if len(oi.atys) != len(i.atys) {
		return false
	}

	for n := range i.atys {
		_, has := oi.atys[n]
		if !has {
			return false
		}
		if !i.atys[n].Same(oi.atys[n]) {
			return false
		}
	}

	return true
}

func (i objectImpl) Equal(a, b Value) Value {
	// The wrapper on type Value guarantees that both values are of the
	// same type, so we can assume they have the same attributes set and
	// just worry about comparing the values for equality.
	for n, ar := range a.v.(map[string]interface{}) {
		br := b.v.(map[string]interface{})[n]
		aty := i.atys[n]

		av := Value{v: ar, ty: aty}
		bv := Value{v: br, ty: aty}

		eq := av.Equal(bv)
		if eq.IsUnknown() {
			return UnknownVal(Bool)
		}
		if !eq.True() {
			return False
		}
	}

	return True
}

func (i objectImpl) ValueSame(a, b Value) bool {
	// This is similar to Equal except that we use the "Same" method to
	// compare the elements, rather than "Equal".

	for n, ar := range a.v.(map[string]interface{}) {
		br := b.v.(map[string]interface{})[n]
		aty := i.atys[n]

		av := Value{v: ar, ty: aty}
		bv := Value{v: br, ty: aty}

		if !av.Same(bv) {
			return false
		}
	}

	return true
}

func init() {
	EmptyObjectVal = Value{
		ty: EmptyObject,
		v:  map[string]interface{}{},
	}
}
