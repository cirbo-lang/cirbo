package cbty

import (
	"testing"

	"github.com/cirbo-lang/cirbo/source"
)

func TestModelType(t *testing.T) {
	ty := Model(testModelImpl{})
	hello := "hello"
	no := "no"
	val1 := ModelVal(ty, &hello)
	val2 := ModelVal(ty, &hello)
	val3 := ModelVal(ty, &no)
	unk := UnknownVal(ty)

	if got, want := ty.Name(), "TestModel"; got != want {
		t.Errorf("wrong Name %#v; want %#v", got, want)
	}

	if got, want := val1.Equal(val2).True(), true; got != want {
		t.Errorf("wrong result for val1.Equal(val2) %#v; want %#v", got, want)
	}
	if got, want := val1.Equal(val3).True(), false; got != want {
		t.Errorf("wrong result for val1.Equal(val2) %#v; want %#v", got, want)
	}

	if got, want := val1.GetAttr("world"), StringVal("hello.world"); !got.Same(want) {
		t.Errorf("wrong result for val1.GetAttr(\"world\") %#v; want %#v", got, want)
	}
	if got, want := unk.GetAttr("world"), UnknownVal(String); !got.Same(want) {
		t.Errorf("wrong result for unk.GetAttr(\"world\") %#v; want %#v", got, want)
	}

	{
		got, _ := val1.Call(CallArgs{})
		want := StringVal("called")
		if !got.Same(want) {
			t.Errorf("wrong result for unk.Call(CallArgs{}) %#v; want %#v", got, want)
		}
	}
	{
		got, _ := unk.Call(CallArgs{})
		want := UnknownVal(String)
		if !got.Same(want) {
			t.Errorf("wrong result for unk.Call(CallArgs{}) %#v; want %#v", got, want)
		}
	}

	if got, want := *(val1.UnwrapModel().(*string)), "hello"; got != want {
		t.Errorf("wrong result for val1.UnwrapModel() %#v; want %#v", got, want)
	}
}

type testModelImpl struct {
}

func (i testModelImpl) Name() string {
	return "TestModel"
}

func (i testModelImpl) SuitableValue(raw interface{}) bool {
	_, isString := raw.(*string)
	return isString
}

func (i testModelImpl) GetAttr(raw interface{}, name string) Value {
	if raw == nil {
		return UnknownVal(String)
	}
	return StringVal(*(raw.(*string)) + "." + name)
}

func (i testModelImpl) CallSignature() *CallSignature {
	return &CallSignature{
		Result: String,
	}
}

func (i testModelImpl) Call(callee interface{}, args CallArgs) (Value, source.Diags) {
	return StringVal("called"), nil
}
