package config

// DefaultTarget returns a target object that represents the machine the binary is currently running on.
func DefaultTarget() *MachineSpec {
	return &MachineSpec{
		Name: "local",
		Machine: map[string]Machine{
			"local": Machine{
				Kind: KindLocal,
			},
		},
	}
}
