package ast

// VisitFunc is the callback signature for VisitAll.
type VisitFunc func(node Node)

// VisitAll is a basic way to traverse the AST beginning with a particular
// node. The given function will be called once for each AST node in
// depth-first order, but no context is provided about the shape of the tree.
func VisitAll(node Node, f VisitFunc) {
	f(node)
	node.walkChildNodes(func(node Node) {
		VisitAll(node, f)
	})
}

// Walker is an interface used with Walk.
type Walker interface {
	// EnterNode is called for each node visited in the traversal.
	//
	// If it returns true, traversal will recurse into child nodes and
	// then eventually ExitNode will be called with the same node. If it
	// returns false, no child nodes are visited and ExitNode is not called.
	EnterNode(node Node) bool

	// ExitNode is called after all of the children of node have been
	// visited.
	ExitNode(node Node)
}

// Walk is a more complex way (than VisitAll) to traverse the AST starting with
// a particular node, which provides information about the tree structure via
// separate Enter and Exit functions.
func Walk(node Node, w Walker) {
	enter := w.EnterNode(node)
	if !enter {
		return
	}
	node.walkChildNodes(func(node Node) {
		Walk(node, w)
	})
	w.ExitNode(node)
}
