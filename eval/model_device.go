package eval

import (
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/source"
)

type device struct {
	name    string
	callSig *cbty.CallSignature
	attrs   StmtBlockAttrs
	block   StmtBlock
}

type deviceModelImpl struct {
	dev *device
}

func deviceValue(dev *device) cbty.Value {
	ty := cbty.Model(deviceModelImpl{dev})
	return cbty.ModelVal(ty, dev)
}

func (i deviceModelImpl) Name() string {
	return "Device"
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
	// TODO: implement
	panic("not callable yet")
}

type deviceInstance struct {
	device *device
	block  *StmtBlockResult
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
	return "Device"
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
