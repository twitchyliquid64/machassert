package config

import (
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
