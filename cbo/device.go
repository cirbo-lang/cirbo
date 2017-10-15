package cbo

// Device represents a reusable device, which can be instantiated multiple
// times in a program.
//
// A device is similar to a circuit except that its internal structure is not
// defined due to being encapsulated into some sort of package, such as
// an integrated circuit.
type Device struct {
	Name      string
	Attrs     AttributesDef
	Terminals TerminalsDef
}

type DeviceInstance struct {
	Device     *Device
	Name       string
	Designator string
	Attrs      map[string]Any
	Terminals  map[string]*TerminalInstance
}
