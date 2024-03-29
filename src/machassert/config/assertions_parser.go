package config

import (
	"errors"
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// ParseAssertionsSchema takes a target configuration and translates it into in-memory structures.
func ParseAssertionsSchema(data []byte) (*AssertionSpec, error) {
	astRoot, err := hcl.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	if _, ok := astRoot.Node.(*ast.ObjectList); !ok {
		return nil, errors.New("schema malformed")
	}

	var outSpec AssertionSpec
	err = hcl.DecodeObject(&outSpec, astRoot)
	if err != nil {
		return nil, err
	}

	normaliseAssertionSpec(&outSpec)
	err = checkAssertionSpec(&outSpec)
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

func normaliseAssertionSpec(spec *AssertionSpec) {
	for i := range spec.Assertions {
		normaliseAssertion(spec.Assertions[i])
		if spec.Assertions[i].Order == 0 {
			spec.Assertions[i].Order = 1000
		}
	}
}

func normaliseAssertion(assertion *Assertion) {
	if len(assertion.Actions) == 0 {
		assertion.Actions = []*Action{&Action{Kind: ActionFail}}
	} else {
		for x := range assertion.Actions {
			if assertion.Actions[x].Kind == "" {
				assertion.Actions[x].Kind = ActionFail
			} else if assertion.Actions[x].Kind == ActionAssert {
				for _, assertion := range assertion.Actions[x].Assertions {
					normaliseAssertion(assertion)
				}
			}
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
	case HashFileAssrt:
		if a.BasePath == "" || a.FilePath == "" {
			return errors.New("base_path/file_path must be specified for file_match assertions")
		}
	case RegexMatchAssrt:
		if a.Regex == "" {
			return errors.New("regex must be specified for regex_contents_match assertions")
		}
	default:
		return errors.New("unsupported assertion type/kind: " + a.Kind)
	}

	for _, action := range a.Actions {
		switch action.Kind {
		case ActionFail:
		case ActionAssert:
			if len(action.Assertions) == 0 {
				return errors.New("at least one assertion must exist for ASSERT actions")
			}
		case ActionCopyFile:
			if action.SourcePath == "" || action.DestinationPath == "" {
				return errors.New("source_path/destination_path must be specified for COPY actions")
			}
		default:
			return errors.New("unsupported action type/kind: " + action.Kind)
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
