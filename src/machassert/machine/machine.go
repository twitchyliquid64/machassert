package machine

// LocalMachine represents the current host as an assertion target.
type LocalMachine struct {
	Name string
}

// SSHMachine represents an assertion target which is communicated with over SSH.
type SSHMachine struct {
	Destination string
	Name        string
}
