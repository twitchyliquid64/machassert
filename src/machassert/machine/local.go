package machine

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"machassert/util"
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

// Run executes the specified command, returning output.
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

// Hash returns the MD5 hash of the file at the given path.
func (m *Local) Hash(fpath string) ([]byte, error) {
	switch runtime.GOOS {
	case "linux":
		o, err := m.Run("md5sum", []string{util.PathSanitize(fpath)})
		if err != nil {
			return nil, err
		}
		return hex.DecodeString(string(o[:32]))
	case "darwin":
		o, err := m.Run("md5", []string{"-q", util.PathSanitize(fpath)})
		if err != nil {
			return nil, err
		}
		hashStr := strings.Trim(string(o), "\n\t ")
		return hex.DecodeString(hashStr)
	default:
		return nil, errors.New("unsupported platform: " + runtime.GOOS)
	}
}

// Grep returns true if the a line in a file match some regular expression.
func (m *Local) Grep(fpath, regex string) (bool, error) {
	_, err := m.Run("grep", []string{"-q", "-E", regex, util.PathSanitize(fpath)})
	if err != nil {
		if _, nonZeroExitStatus := err.(*exec.ExitError); nonZeroExitStatus {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ReadFile returns a reader to a file on a local machine.
func (m *Local) ReadFile(fpath string) (io.ReadCloser, error) {
	return os.Open(util.PathSanitize(fpath))
}

// Close releases the resources associated with the machine.
func (m *Local) Close() error {
	return nil
}

// WriteFile returns a writer which can be used to write content to the remote file.
func (m *Local) WriteFile(fpath string) (io.WriteCloser, error) {
	return os.OpenFile(util.PathSanitize(fpath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
}
