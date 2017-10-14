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
