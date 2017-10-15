package eval

import (
	"github.com/cirbo-lang/cirbo/cbo"
	"github.com/cirbo-lang/cirbo/cbty"
)

// Unwrapper is a type that can be used to "unwrap" cbty.Value instances to
// more convenient native types, via the intermediate interface type cbo.Any.
//
// An Unwrapper instance serves as a cache for unwrapped values, so that
// unwrapping the same object twice can (in most cases, at least) produce
// the same object.
//
// The zero value of Unwrapper is an unwrapper ready to use, though most
// callers will want to create a pointer to that zero value.
type Unwrapper struct {
	devices         map[*device]*cbo.Device
	deviceInstances map[*deviceInstance]*cbo.DeviceInstance
}

// Unwrap obtains a native Go value corresponding to the given value within
// Cirbo's type system. See the documentation of type cbo.Any for details
// on what can be returned by this function.
//
// This function is not concurrency-safe. If multiple goroutines are sharing
// an Unwrapper then the caller must use a mutex or other similar locking
// primitive to prevent concurrent calls to Unwrap.
func (u *Unwrapper) Unwrap(val cbty.Value) cbo.Any {
	if val == cbty.NilValue || val.IsUnknown() {
		return nil
	}

	ty := val.Type()

	// Simple types first
	switch ty {
	case cbty.String:
		return val.AsString()
	case cbty.Bool:
		return val.True()
	case cbty.TypeType:
		return val.UnwrapType()
	}

	switch {
	case ty.IsModel():
		return u.unwrapModel(val.UnwrapModel())
	case ty.IsNumber():
		return val.AsQuantity()
	}

	// FIXME: should have a mapping for every possible cbty type
	return nil
}

func (u *Unwrapper) unwrapModel(raw interface{}) cbo.Any {
	switch tv := raw.(type) {
	case *device:
		if u.devices != nil && u.devices[tv] != nil {
			return u.devices[tv]
		}
		ret := tv.AsPublic()
		if u.devices == nil {
			u.devices = map[*device]*cbo.Device{}
		}
		u.devices[tv] = ret
		return ret
	case *deviceInstance:
		if u.deviceInstances != nil && u.deviceInstances[tv] != nil {
			return u.deviceInstances[tv]
		}
		ret := &cbo.DeviceInstance{
			Name:       tv.name,
			Device:     u.unwrapModel(tv.device).(*cbo.Device),
			Designator: tv.content.Designator,
			Attrs:      map[string]cbo.Any{},
		}
		if ret.Designator == "" {
			ret.Designator = "X"
		}
		for name, attr := range tv.device.attrs {
			val := tv.content.Context.Value(attr.Symbol)
			ret.Attrs[name] = u.Unwrap(val)
		}
		if u.deviceInstances == nil {
			u.deviceInstances = map[*deviceInstance]*cbo.DeviceInstance{}
		}
		u.deviceInstances[tv] = ret
		return ret
	default:
		// Should never happen, since we should exhaustively cover
		// all of our model types in here.
		return nil
	}
}
