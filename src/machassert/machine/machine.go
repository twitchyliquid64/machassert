package machine

import "machassert/config"

type authPromptProvider interface {
	AuthenticationPrompt(prompt string) (string, error)
}

// SSH represents an assertion target which is communicated with over SSH.
type SSH struct {
	Destination string
	Name        string
}

// ConnectLocal returs a Local machine object
func ConnectLocal(name string, machine *config.Machine) (*Local, error) {
	if machine.Kind != config.KindLocal {
		panic("machine kind must be local")
	}
	return &Local{
		MachineName: name,
	}, nil
}
