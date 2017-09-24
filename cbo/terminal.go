package cbo

type TerminalsDef struct {
	All   map[string]Terminal
	Names []string
}

func (td TerminalsDef) Ordered() []Terminal {
	ret := make([]Terminal, len(td.Names))
	for i, n := range td.Names {
		ret[i] = td.All[n]
	}
	return ret
}

type Terminal struct {
	Name string

	// LowerBound and UpperBound are both zero for single-connection terminals,
	// and will have other values for bus-type terminals. UpperBound must
	// always be greater than LowerBound when both values are not zero.
	LowerBound int
	UpperBound int

	Type       TerminalType
	Dir        TerminalDir
	Role       TerminalRole
	OutputType TerminalOutputType
}

func (t *Terminal) NewInstance() *TerminalInstance {
	endpointCt := (t.UpperBound - t.LowerBound) + 1
	outside := make([]*Endpoint, endpointCt)
	inside := make([]*Endpoint, endpointCt)
	for i := range outside {
		outside[i] = &Endpoint{
			Name:       t.Name,
			Net:        nil, // none yet
			Type:       t.Type,
			Dir:        t.Dir,
			Role:       t.Role,
			OutputType: t.OutputType,
		}
		inside[i] = &Endpoint{
			Name: t.Name,
			Net:  nil, // none yet
			Type: t.Type,
			Dir:  t.Dir.Inverse(),
			Role: t.Role.Inverse(),

			// If the terminal is an input then we can't automatically infer
			// the output type of its "inside" because that depends on what
			// ultimately gets connected to our "outside", so we'll leave it
			// unspecified for now and resolve this once we're "flattening"
			// the overall design into a single circuit.
			OutputType: NoOutput,
		}
	}
	return &TerminalInstance{
		Terminal: t,
		Outside:  outside,
		Inside:   inside,
	}
}

type TerminalInstance struct {
	Terminal *Terminal
	Outside  []*Endpoint
	Inside   []*Endpoint
}

type TerminalType rune

//go:generate stringer -type=TerminalType

const (
	Passive TerminalType = 0
	Signal  TerminalType = 'S'
	Power   TerminalType = 'P'
)

type TerminalDir rune

//go:generate stringer -type=TerminalDir

const (
	Undirected TerminalDir = 0
	Input      TerminalDir = 'I'
	Output     TerminalDir = 'O'

	// Bidirectional is a terminal direction that indicates that a terminal
	// operates both as an input and and output.
	//
	// Bidirectional terminals must always have their role set to either
	// Leader or Follower. NoRole is not a permitted role for bidirectional
	// terminals.
	Bidirectional TerminalDir = 'B'
)

type TerminalOutputType rune

//go:generate stringer -type=TerminalOutputType

const (
	NoOutput      TerminalOutputType = 0
	PushPull      TerminalOutputType = 'P'
	Tristate      TerminalOutputType = '±'
	OpenCollector TerminalOutputType = 'C'
	OpenEmitter   TerminalOutputType = 'E'
)

// TerminalRole is an enumeration describing different roles a terminal can
// play from the standpoint of logical system structure.
//
// The roles on the terminals of devices are used to imply roles of devices
// themselves in relation to one another. For example, in many systems there
// is a single microcontroller behaving as the logical "leader", directing
// the operation of one or more other devices that are "followers" in this
// sense.
//
// The relationships between devices implied by their roles are intended to
// communicate a similar idea as the horizontal axis of a good schematic
// diagram: a device that "leads" would generally be presented to the left
// of a device that "follows", so that the delegation of responsibility
// can be easily understood via visual relationships.
//
// In more complex systems, a particular device (or sub-circuit) may play
// multiple roles. For example, an SPI-driven Ethernet controller could be
// considered to be a follower on its SPI interface (where another device
// instructs it to send frames) but a leader on its Ethernet interface.
// This is why it is terminals, rather than devices themselves, that have
// roles specified.
type TerminalRole rune

//go:generate stringer -type=TerminalRole

const (
	NoRole TerminalRole = 0

	// Leader is a terminal role that indicates that the object to which
	// the terminal belongs is the initiator of communication.
	//
	// For example, the MISO input on a microcontroller's SPI bus should
	// generally be an "input leader" since it is the microcontroller that
	// is in control of the conversation (via the associated SCLK signal).
	//
	// As an exception, interrupt pins (for asynchronous notification
	// "upstream" from a following device to a leading device) should
	// be "follower" on the _sending_ end, even though that device is the
	// initiator, to preserve the leader-follower relationship between
	// the devices themselves.
	Leader TerminalRole = 'L'

	// Follower is a terminal role that indicates that the object to which
	// the terminal belongs is a non-initiating participant in
	// communication.
	//
	// For example, the MISO output on a SPI "slave" device should generally
	// be an "output follower" since it only produces signals in response
	// to requests from the SPI "master" (via the associated SCLK signal).
	//
	// If a device is able to switch roles between leader and follower,
	// Leader should take precedence in deciding a role.
	Follower TerminalRole = 'F'
)

func (d TerminalDir) Inverse() TerminalDir {
	switch d {
	case Input:
		return Output
	case Output:
		return Input
	case Bidirectional:
		return Bidirectional
	default:
		return Undirected
	}
}

func (r TerminalRole) Inverse() TerminalRole {
	switch r {
	case Leader:
		return Follower
	case Follower:
		return Leader
	default:
		return NoRole
	}
}
