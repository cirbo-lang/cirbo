package cbo

import (
	"github.com/cirbo-lang/cirbo/cty"
)

// Device represents a reusable device, which can be instantiated multiple
// times in a program.
//
// A device is similar to a circuit except that its internal structure is not
// defined due to being encapsulated into some sort of package, such as
// an integrated circuit.
type Device struct {
	Name       string
	Designator string
	Attrs      AttributesDef
	Terminals  TerminalsDef
}

func (d *Device) NewInstance(name string, attrs map[string]cty.Value) *DeviceInstance {
	// Like all terminals, these have both inside and outside endpoints
	// but we will use only the outside ones on a device since its
	// internals are opaque.
	terminals := make(map[string]*TerminalInstance, len(d.Terminals.All))
	for name, terminal := range d.Terminals.All {
		terminals[name] = terminal.NewInstance()
	}

	return &DeviceInstance{
		Device:    d,
		Name:      name,
		Attrs:     attrs,
		Terminals: terminals,
	}
}

type DeviceInstance struct {
	Device    *Device
	Name      string
	Attrs     map[string]cty.Value
	Terminals map[string]*TerminalInstance
}

// DevicesDef describes the devices instantiated in a particular context.
//
// The keys of this map are the symbol names to be used for device instances
// within that context, which are not (necessarily) the same as the device
// names themselves, which are global.
type DevicesDef map[string]Device
