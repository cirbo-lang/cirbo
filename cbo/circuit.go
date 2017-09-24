package cbo

import (
	"github.com/cirbo-lang/cirbo/cty"
)

// Circuit represents a reusable circuit, which can be instantiated multiple
// times in a program.
//
// Instances of a Circuit are represented by CircuitInstance.
type Circuit struct {
	Name      string
	Attrs     AttributesDef
	Terminals TerminalsDef
}

func (c *Circuit) NewInstance(name string, attrs map[string]cty.Value) *CircuitInstance {
	// We need to recursively instantiate all of the contained items

	terminals := make(map[string]*TerminalInstance, len(c.Terminals.All))
	for name, terminal := range c.Terminals.All {
		terminals[name] = terminal.NewInstance()
	}

	// TODO: Devices
	// TODO: Circuits

	return &CircuitInstance{
		Circuit:   c,
		Name:      name,
		Attrs:     attrs,
		Terminals: terminals,
	}
}

// CircuitInstance represents a particular concrete instantiation of a
// Circuit, which provides a set of values for the circuit's attributes and
// describes how its terminals are connected to other devices and circuits
// within the same heirarchy level.
type CircuitInstance struct {
	Circuit   *Circuit
	Name      string
	Attrs     map[string]cty.Value
	Terminals map[string]*TerminalInstance
	Devices   map[string]*DeviceInstance
	Circuits  map[string]*CircuitInstance
}

// CircuitsDef represents the circuits instantiated in a particular context.
//
// The keys of this map are the symbol names to be used for circuit instances
// within that context, which are not (necessarily) the same as the circuit
// names themselves, which are global.
type CircuitsDef map[string]Circuit
