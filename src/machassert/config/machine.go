package config

const (
	// KindLocal represents the machine the binary is executing on.
	KindLocal string = "local"
	// KindSSH represents a machine accessible via SSH.
	KindSSH string = "ssh"

	//AuthKindPassword represents password authentication
	AuthKindPassword = "password"
)

//MachineSpec describes the high-level schema for target configuration.
type MachineSpec struct {
	Name    string
	Machine map[string]Machine
}

//Machine describes the target schema for a specific machine.
type Machine struct {
	Kind        string
	Destination string //only valid for non local machines
	Auth        []MachineAuth
}

//MachineAuth describes the scheme for machine authentication configuration.
type MachineAuth struct {
	Kind     string
	Password string
}
