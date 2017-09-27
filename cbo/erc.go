package cbo

// ERCMode represents the characteristics of an endpoint that are used for
// electrical rules check.
type ERCMode struct {
	Type       ERCType
	Dir        ERCDir
	OutputType ERCOutputType
}

type ERCType rune

//go:generate stringer -type=ERCType

const (
	Passive ERCType = 0
	Signal  ERCType = 'S'
	Power   ERCType = 'P'
)

type ERCDir rune

//go:generate stringer -type=ERCDir

const (
	Undirected ERCDir = 0
	Input      ERCDir = 'I'
	Output     ERCDir = 'O'

	// Bidirectional is a terminal direction that indicates that a terminal
	// operates both as an input and and output.
	//
	// Bidirectional terminals must always have their role set to either
	// Leader or Follower. NoRole is not a permitted role for bidirectional
	// terminals.
	Bidirectional ERCDir = 'B'

	// ERCDir is also overloaded to communicate a number of special terminal
	// types that act as flags to customize behavior in the rules checker.
	// These cannot be set directly in the language, but some standard library
	// functions (implemented in Go) produce terminals with these directions
	// to help the user explain unusual situations to the checker.

	// MultiOutputSinkFlag is a special "direction" that represents an endpoint
	// that can be driven by multiple outputs, where that would usually be
	// an error. A net with an endpoint of this direction must otherwise have
	// only outputs.
	MultiOutputSinkFlag ERCDir = 'Ⓜ'

	// NoConnectFlag is a special "direction" that represents that another
	// endpoint is intentionally not connected, where that would usually be an
	// error. A net with an endpoint of this direction must have exactly two
	// endpoints, one of which is this flag and the other is the endpoint that
	// is intentionally not connected.
	NoConnectFlag ERCDir = 'Ⓝ'
)

type ERCOutputType rune

//go:generate stringer -type=ERCOutputType

const (
	NoOutput      ERCOutputType = 0
	PushPull      ERCOutputType = 'P'
	Tristate      ERCOutputType = '±'
	OpenCollector ERCOutputType = 'C'
	OpenEmitter   ERCOutputType = 'E'
	UnknownOutput ERCOutputType = '?'
)

// Inverse returns the ERC mode for the "other side" of a particular ERC mode.
// That is, and input becomes and output and vice-versa.
//
// This is used to create the ERC mode for the "inside" of a terminal, since
// the user describes a terminal from the perspective of an outside caller.
//
// Inverse is a lossy operation: inverting an input produces an output with
// an unknown output type, because the output type depends on what the
// other side ends up being connected to. Conversely, inverting an output
// to produce an input loses its output type.
func (d ERCMode) Inverse() ERCMode {
	var oType ERCOutputType
	oDir := d.Dir.Inverse()
	if oDir == Output || oDir == Bidirectional {
		oType = UnknownOutput
	}
	return ERCMode{
		Type:       d.Type,
		Dir:        d.Dir.Inverse(),
		OutputType: oType,
	}
}

func (d ERCDir) Inverse() ERCDir {
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
