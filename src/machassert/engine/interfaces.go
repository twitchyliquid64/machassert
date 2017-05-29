package engine

import (
	"io"
)

// Machine represents a target for assertions. The base type implements the communication layer to the target.
type Machine interface {
	Name() string
	ReadFile(fpath string) (io.ReadCloser, error)
}
