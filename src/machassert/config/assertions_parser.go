package config

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// ParseAssertionsSchema takes a target configuration and translates it into in-memory structures.
func ParseAssertionsSchema(data []byte) (*AssertionSpec, error) {
	astRoot, err := hcl.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	var outSpec AssertionSpec
	err = hcl.DecodeObject(&outSpec, astRoot)
	if err != nil {
		return nil, err
	}

	return &outSpec, nil
}

// ParseAssertionsSpecFile parses the assertions file from disk.
func ParseAssertionsSpecFile(fpath string) (*AssertionSpec, error) {
	d, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ParseAssertionsSchema(d)
}
