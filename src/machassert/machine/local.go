package machine

import (
	"io"
	"os"
)

// Local represents the current host as an assertion target.
type Local struct {
	MachineName string
}

// Name returns the name of the target
func (m *Local) Name() string {
	return m.MachineName
}

// ReadFile returns a reader to a file on a local machine.
func (m *Local) ReadFile(fpath string) (io.ReadCloser, error) {
	return os.Open(fpath)
}
