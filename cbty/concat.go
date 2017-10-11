package cbty

// typeWithConcat is an interface implemented by typeImpls that can
// concatenate.
type typeWithConcat interface {
	CanConcat(other Type) bool
	Concat(a, b Value) Value
}
