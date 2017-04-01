package engine

// Machine represents a target for assertions. The base type implements the communication layer to the target.
type Machine interface {
	Name() string
}
