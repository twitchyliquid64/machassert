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

func checkAssertion(a *Assertion) error {
	switch a.Kind {
	case FileExistsAssrt:
		fallthrough
	case FileNotExistsAssrt:
		if a.FilePath == "" {
			return errors.New("file_path must be specified for exists and !exists assertions")
		}
	case HashMatchAssrt:
		if a.Hash == "" || a.FilePath == "" {
			return errors.New("hash/file_path must be specified for md5_match assertions")
		}
	default:
		return errors.New("unsupported assertion type/kind")
	}

	for _, action := range a.Actions {
		switch action.Kind {
		case ActionFail:
		case ActionApplyFile:
			if action.SourcePath == "" || action.DestinationPath == "" {
				return errors.New("source_path/destination_path must be specified for APPLY actions")
			}
		default:
			return errors.New("unsupported action type/kind")
		}
	}
	return nil
}

func checkAssertionSpec(spec *AssertionSpec) error {
	for name, a := range spec.Assertions {
		if name == "" {
			return errors.New("name must be specified for an assertion")
		}
		err := checkAssertion(a)
		if err != nil {
			return err
		}
	}
	return nil
}
