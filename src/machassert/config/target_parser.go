package config

import (
	"errors"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// ParseTargetSchema takes a target configuration and translates it into in-memory structures.
func ParseTargetSchema(data []byte) (*MachineSpec, error) {
	astRoot, err := hcl.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	var outSpec MachineSpec
	err = hcl.DecodeObject(&outSpec, astRoot)
	if err != nil {
		return nil, err
	}

	err = normalizeMachineSpec(&outSpec)
	if err != nil {
		return nil, err
	}
	err = validateMachineSpec(&outSpec)
	if err != nil {
		return nil, err
	}

	return &outSpec, nil
}

// ParseTargetSpecFile parses the targets file from disk.
func ParseTargetSpecFile(fpath string) (*MachineSpec, error) {
	d, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ParseTargetSchema(d)
}

func normalizeMachineSpec(spec *MachineSpec) error {
	for k := range spec.Machine {
		for i := range spec.Machine[k].Auth {
			if spec.Machine[k].Auth[i].Password != "" { //If password is set, set the auth kind to password
				spec.Machine[k].Auth[i].Kind = AuthKindPassword
			}
		}
	}
	return nil
}

func validateMachineSpec(spec *MachineSpec) error {
	for k := range spec.Machine {
		for i := range spec.Machine[k].Auth {
			switch spec.Machine[k].Auth[i].Kind {
			case AuthKindPassword:
				if spec.Machine[k].Auth[i].Password == "" {
					return errors.New("Password must be specified for password authentication")
				}
			default:
				return errors.New("Invalid machine auth type")
			}
		}
	}
	return nil
}
