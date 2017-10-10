package ast

// SwapTable is a mapping table from source nodes to replacement nodes, used
// to represent on-the-fly AST transformations without necessarily modifying
// the original AST.
type SwapTable map[Node]Node

// Add inserts a new entry into the table.
//
// This method will panic on a nil SwapTable.
func (t SwapTable) Add(match, replacement Node) {
	t[match] = replacement
}

// Swap takes a node and returns either its replacement (if it has one) or
// the original node.
//
// Calling Swap on a nil SwapTable is allowed, and will behave as if called
// on an empty table.
func (t SwapTable) Swap(node Node) Node {
	if t == nil {
		return node
	}
	if new, swapped := t[node]; swapped {
		return new
	}
	return node
}
