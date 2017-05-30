package config

import (
	"errors"
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

	err = checkAssertionSpec(&outSpec)
	if err != nil {
		return nil, err
	}

	normaliseAssertionSpec(&outSpec)
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

func normaliseAssertionSpec(spec *AssertionSpec) {
	for i := range spec.Assertions {
		if len(spec.Assertions[i].Actions) == 0 {
			spec.Assertions[i].Actions = []*Action{&Action{Kind: ActionFail}}
		}
	}
}

func checkAssertionSpec(spec *AssertionSpec) error {
	for name, a := range spec.Assertions {
		if name == "" {
			return errors.New("name must be specified for an assertion")
		}
		switch a.Kind {
		case FileExistsAssrt:
		case FileNotExistsAssrt:
		default:
			return errors.New("unsupported assertion type/kind")
		}
	}
	return nil
}
