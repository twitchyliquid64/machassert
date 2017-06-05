package machine

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Local represents the current host as an assertion target.
type Local struct {
	MachineName string
}

// Name returns the name of the target
func (m *Local) Name() string {
	return m.MachineName
}

func (m *Local) Run(name string, args []string) ([]byte, error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// MD5 returns the hash of the file at the given path.
func (m *Local) Hash(fpath string) ([]byte, error) {
	switch runtime.GOOS {
	case "darwin":
		o, err := m.Run("md5", []string{"-q", fpath})
		if err != nil {
			return nil, err
		}
		hashStr := strings.Trim(string(o), "\n\t ")
		return hex.DecodeString(hashStr)
	default:
		return nil, errors.New("unsupported platform: " + runtime.GOOS)
	}
}

// ReadFile returns a reader to a file on a local machine.
func (m *Local) ReadFile(fpath string) (io.ReadCloser, error) {
	return os.Open(fpath)
}

func (m *Local) Close() error {
	return nil
}
