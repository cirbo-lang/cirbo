package eval

import (
	"github.com/cirbo-lang/cirbo/source"
)

// range can be embedded into an expression to implement the sourceRange method.
type rng struct {
	r source.Range
}

// srcRange is a convenience function for concisely populating rng instances.
func srcRange(r source.Range) rng {
	return rng{
		r: r,
	}
}

func (r rng) sourceRange() source.Range {
	return r.r
}

// leafExpr can be embedded into an expression type that has no child
// expressions, to get a do-nothing eachChild implementation.
type leafExpr struct {
}

func (leafExpr) eachChild(cb walkCb) {
	// nothing to do!
}
