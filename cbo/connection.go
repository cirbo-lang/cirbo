package cbo

// Connection represents a connection between two objects in a particular
// context.
//
// Connections are used within circuits to represent the internal connectivity
// that needs to be re-created for each instance of the circuit.
type Connection struct {
}

// ConnectionEndpoint represents one end of a connection.
type ConnectionEndpoint struct {
	ObjectName string
}
