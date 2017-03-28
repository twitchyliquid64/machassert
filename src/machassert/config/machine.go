package config

type machineKind string

const (
	// KindLocal represents the machine the binary is executing on.
	KindLocal machineKind = "local"
	// KindSSH represents a machine accessible via SSH.
	KindSSH machineKind = "ssh"
)

//MachineSpec describes the high-level schema for target configuration.
type MachineSpec struct {
	Name     string
	Machines []Machine
}

//Machine describes the target schema for a specific machine.
type Machine struct {
	Kind        machineKind
	Destination string //only valid for non local machines
}
