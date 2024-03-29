package config

// Valid machine types
const (
	KindLocal string = "local"
	KindSSH   string = "ssh"
)

// Valid authentication means/types
const (
	AuthKindPassword = "password"
	AuthKindPrompt   = "prompt"
	AuthKindLocalKey = "user-key"
	AuthKindKeyFile  = "key-file"
)

// MachineSpec describes the high-level schema for target configuration.
type MachineSpec struct {
	Name    string
	Machine map[string]*Machine
}

//Machine describes the target schema for a specific machine.
type Machine struct {
	Kind        string
	Destination string //only valid for non local machines
	Username    string //only needed for SSH
	Auth        []MachineAuth
}

//MachineAuth describes the scheme for machine authentication configuration.
type MachineAuth struct {
	Kind     string
	Password string
	Key      string
}
