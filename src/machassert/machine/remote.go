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

type Remote struct {
	MachineName string
	Address     string
	authInfo    []config.MachineAuth
	conn        *ssh.Client
}

func ConnectRemote(name string, m *config.Machine) (*Remote, error) {
	c := &ssh.ClientConfig{User: m.Username, HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	for _, authItem := range m.Auth {
		switch authItem.Kind {
		case config.AuthKindPassword:
			c.Auth = append(c.Auth, ssh.Password(authItem.Password))
		}
	}
	//TODO: Add ability to verify host key

	client, err := ssh.Dial("tcp", m.Destination, c)
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

// MD5 returns the hash of the file at the given path.
func (m *Remote) Hash(fpath string) ([]byte, error) {
	o, err := m.Run("md5sum", []string{fpath})
	if err != nil {
		return nil, err
	}
	hashStr := strings.Trim(strings.Split(string(o), " ")[0], "\n\t ")
	return hex.DecodeString(hashStr)
}

func (r *Remote) Name() string {
	return r.MachineName
}

func (r *Remote) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}
