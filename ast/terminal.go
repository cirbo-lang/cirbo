package ast

type Terminal struct {
	WithRange
	Name       string
	Type       TerminalType
	Dir        TerminalDir
	OutputType TerminalOutputType
}

type TerminalType rune

const (
	Passive TerminalType = 0
	Signal  TerminalType = 'S'
	Power   TerminalType = 'P'
)

type TerminalDir rune

const (
	Undirected TerminalDir = 0
	Input      TerminalDir = 'I'
	Output     TerminalDir = 'O'

	// BidiLeader is a terminal direction that indicates that it operates
	// both as an input and an output but that the associated device is
	// an "bus leader", meaning that it initiates transactions.
	//
	// For example, in a typical application with a microcontroller requesting
	// data from a peripheral using I2C, the microcontroller would be a
	// BidiLeader and the peripheral a BidiFollower, since the peripheral
	// remains inactive until the microcontroller requests data.
	BidiLeader TerminalDir = 'L'

	// BidiFollower is a terminal direction that indicates that it operates
	// both as an input and an output but that the associated device is
	// a "bus follower", meaning that it participates in transactions
	// initiated by other devices on the bus but does not itself begin
	// transactions.
	//
	// If a device participates both as a leader and a follower, BidiLeader
	// should be considered to take precedence.
	BidiFollower TerminalDir = 'F'
)

type TerminalOutputType rune

const (
	NoOutput      TerminalOutputType = 0
	PushPull      TerminalOutputType = 'P'
	Tristate      TerminalOutputType = '±'
	OpenCollector TerminalOutputType = 'C'
	OpenEmitter   TerminalOutputType = 'E'
)

func (n *Terminal) walkChildNodes(cb internalWalkFunc) {
	// Terminal is a leaf node
}
