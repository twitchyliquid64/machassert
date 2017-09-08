package machine

import (
	"bytes"
	"encoding/hex"
	"io"
	"machassert/config"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Remote is a target connected to via SSH.
type Remote struct {
	MachineName string
	Address     string
	authInfo    []config.MachineAuth
	conn        *ssh.Client
}

// ConnectRemote opens an SSH connection to a remote target.
func ConnectRemote(name string, m *config.Machine) (*Remote, error) {
	c := &ssh.ClientConfig{User: m.Username, HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	for _, authItem := range m.Auth {
		switch authItem.Kind {
		case config.AuthKindPassword:
			c.Auth = append(c.Auth, ssh.Password(authItem.Password))
		}
	}

	address := m.Destination
	if !strings.Contains(address, ":") {
		address += ":22"
	}

	//TODO: Add ability to verify host key
	client, err := ssh.Dial("tcp", address, c)
	if err != nil {
		return nil, err
	}

	return &Remote{
		MachineName: name,
		Address:     m.Destination,
		authInfo:    m.Auth,
		conn:        client,
	}, nil
}

type readFileRemoteReadCloser struct {
	buff    bytes.Buffer
	session *ssh.Session
}

func (r *readFileRemoteReadCloser) Read(p []byte) (int, error) {
	return r.buff.Read(p)
}
func (r *readFileRemoteReadCloser) Close() error {
	return r.session.Close()
}

// ReadFile returns a reader to a file on a local machine.
func (r *Remote) ReadFile(fpath string) (io.ReadCloser, error) {
	s, err := r.conn.NewSession()
	if err != nil {
		return nil, err
	}

	err = s.Run("cat " + fpath)
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			if exitError.ExitStatus() == 1 {
				return nil, os.ErrNotExist
			}
		}
		return nil, err
	}

	out := &readFileRemoteReadCloser{session: s}
	s.Stdout = &out.buff
	return out, nil
}

// Run executes the specified command, returning output.
func (r *Remote) Run(name string, args []string) ([]byte, error) {
	var out bytes.Buffer
	s, err := r.conn.NewSession()
	if err != nil {
		return nil, err
	}
	s.Stdout = &out

	for i := range args {
		if strings.ContainsAny(args[i], " |\"'") {
			args[i] = "\"" + strings.Replace(args[i], "\"", "\\\"", -1) + "\""
		}
	}

	err = s.Run(name + " " + strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// Hash returns the MD5 hash of the file at the given path.
func (r *Remote) Hash(fpath string) ([]byte, error) {
	o, err := r.Run("md5sum", []string{fpath})
	if err != nil {
		return nil, err
	}
	hashStr := strings.Trim(strings.Split(string(o), " ")[0], "\n\t ")
	return hex.DecodeString(hashStr)
}

// Name returns the name of the target
func (r *Remote) Name() string {
	return r.MachineName
}

// Close releases the resources associated with the machine.
func (r *Remote) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

type sshWriter struct {
	s *ssh.Session
	w io.WriteCloser
}

func (w *sshWriter) Close() error {
	w.w.Close()
	return w.s.Close()
}

func (w *sshWriter) Write(in []byte) (int, error) {
	return w.w.Write(in)
}

// WriteFile returns a writer which can be used to write content to the remote file.
func (r *Remote) WriteFile(fpath string) (io.WriteCloser, error) {
	s, err := r.conn.NewSession()
	if err != nil {
		return nil, err
	}

	pipe, err := s.StdinPipe()
	if err != nil {
		s.Close()
		return nil, err
	}

	args := []string{"cat - > " + strings.Replace(fpath, "\"", "\\\"", -1)}
	err = s.Start(strings.Join(args, " "))
	if err != nil {
		s.Close()
		return nil, err
	}
	return &sshWriter{s, pipe}, nil
}
