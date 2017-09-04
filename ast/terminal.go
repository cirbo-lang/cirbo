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
	Bidirectional TerminalDir = 0
	Input         TerminalDir = 'I'
	Output        TerminalDir = 'O'
)

type TerminalOutputType rune

const (
	NoOutput      TerminalOutputType = 0
	Tristate      TerminalOutputType = '±'
	OpenCollector TerminalOutputType = 'C'
	OpenDrain     TerminalOutputType = 'D'
)

func (n *Terminal) walkChildNodes(cb internalWalkFunc) {
	// Terminal is a leaf node
}
