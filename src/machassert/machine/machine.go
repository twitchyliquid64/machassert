package machine

// Machine represents a target for assertions. The base type implements the communication layer to the target.
type Machine interface {
	Name() string
}

// SSHMachine represents an assertion target which is communicated with over SSH.
type SSHMachine struct {
	Destination string
	Name        string
}
