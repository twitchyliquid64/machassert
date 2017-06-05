package engine

import (
	"encoding/hex"
	"errors"
	"machassert/config"
	"os"
	"strings"
)

// ErrAssertionsFailed is returned if spec execution is short-circuited due to assertion failure.
var ErrAssertionsFailed = errors.New("assertions failed")

// Assertion Results
const (
	AssertionNoop int = iota
	AssertionApplied
	AssertionFailed
	AssertionError
)

// AssertionResult captures what happens when an assertion is applied.
type AssertionResult struct {
	Result int
}

func (r AssertionResult) String() string {
	switch r.Result {
	case AssertionNoop:
		return "OK"
	case AssertionFailed:
		return "FAILED"
	case AssertionApplied:
		return "APPLIED"
	case AssertionError:
		return "ERR"
	default:
		return "?"
	}
}

func applyAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	result := &AssertionResult{Result: AssertionError}
	var err error

	switch assertion.Kind {
	case config.HashMatchAssrt:
		result, err = applyHashAssertion(machine, assertion)
	case config.FileExistsAssrt:
		result, err = applyExistsAssertion(machine, assertion)
	}

	if err == nil && result.Result == AssertionApplied {
		for _, action := range assertion.Actions {
			err = doAction(machine, assertion, action)
			if err == ErrAssertionsFailed {
				result.Result = AssertionFailed
			}
			if err != nil {
				return result, err
			}
		}
	}

	return result, err
}

func doAction(machine Machine, assertion *config.Assertion, action *config.Action) error {
	switch action.Kind {
	case "":
		return nil
	case "FAIL":
		return ErrAssertionsFailed
	default:
		return errors.New("Unrecognised actions kind")
	}
}

func applyExistsAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	f, err := machine.ReadFile(assertion.FilePath)
	if err != nil && os.IsNotExist(err) {
		return &AssertionResult{Result: AssertionApplied}, nil
	}
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	f.Close()
	return &AssertionResult{Result: AssertionNoop}, nil
}

func applyHashAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	hash, err := machine.Hash(assertion.FilePath)
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	if hex.EncodeToString(hash) != strings.ToLower(assertion.Hash) {
		return &AssertionResult{Result: AssertionApplied}, nil
	}
	return &AssertionResult{Result: AssertionNoop}, nil
}
