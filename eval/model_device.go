package eval

import (
	"fmt"

	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

type device struct {
	name    string
	callSig *cbty.CallSignature
	attrs   StmtBlockAttrs
	block   StmtBlock
	instTy  cbty.Type
}

type deviceModelImpl struct {
	dev *device
}

func deviceValue(dev *device) cbty.Value {
	ty := cbty.Model(deviceModelImpl{dev})
	return cbty.ModelVal(ty, dev)
}

func (i deviceModelImpl) Name() string {
	return fmt.Sprintf("Device(%q)", i.dev.name)
}

func (i deviceModelImpl) SuitableValue(raw interface{}) bool {
	_, isDevice := raw.(*device)
	return isDevice
}

func (i deviceModelImpl) GetAttr(raw interface{}, name string) cbty.Value {
	return cbty.NilValue
}

func (i deviceModelImpl) CallSignature() *cbty.CallSignature {
	return i.dev.callSig
}

func (i deviceModelImpl) Call(callee interface{}, args cbty.CallArgs) (cbty.Value, source.Diags) {
	dev := callee.(*device)
	context := args.Context.(*Context)

	if args.TargetName == "" {
		return cbty.UnknownVal(dev.instTy), source.Diags{
			{
				Level:   source.Error,
				Summary: "Anonymous device instance",
				Detail:  "A device instance may be created only when assigning directly to a name, using an assignment statement.",
				Ranges:  args.CallRange.List(),
			},
		}
	}

	initDefs := make(map[*Symbol]cbty.Value, len(dev.attrs))
	for name, attr := range dev.attrs {
		sym := attr.Symbol
		val, defined := args.Explicit[name]
		if defined {
			initDefs[sym] = val
		} else {
			if attr.Default == cbty.NilValue {
				// Should never happen, but we'll put something safe here
				// anyway so that we won't crash later trying to eval this.
				initDefs[sym] = cbty.PlaceholderVal
				continue
			}
			initDefs[sym] = attr.Default
		}
	}

	result, diags := dev.block.Execute(StmtBlockExecute{
		Context: context,
	}, initDefs)

	inst := &deviceInstance{
		name:    args.TargetName,
		device:  dev,
		content: result,
	}

	return deviceInstanceValue(dev.instTy, inst), diags
}

type deviceInstance struct {
	name    string
	device  *device
	content *StmtBlockResult
}

type deviceInstanceModelImpl struct {
	device *device
}

func deviceInstanceType(dev *device) cbty.Type {
	return cbty.Model(deviceInstanceModelImpl{dev})
}

func deviceInstanceValue(ty cbty.Type, di *deviceInstance) cbty.Value {
	return cbty.ModelVal(ty, di)
}

func (i deviceInstanceModelImpl) Name() string {
	return i.device.name
}

func (i deviceInstanceModelImpl) SuitableValue(raw interface{}) bool {
	di, isInstance := raw.(*deviceInstance)
	if !isInstance {
		return false
	}
	return di.device == i.device
}

func (i deviceInstanceModelImpl) GetAttr(raw interface{}, name string) cbty.Value {
	return cbty.NilValue
}

func (i deviceInstanceModelImpl) CallSignature() *cbty.CallSignature {
	return nil
}

func (i deviceInstanceModelImpl) Call(callee interface{}, args cbty.CallArgs) (cbty.Value, source.Diags) {
	panic("not callable") // should never get here because CallSignature returns nil
}
