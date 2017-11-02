package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"machassert/config"
	"machassert/util"
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
	AssertionApplyError
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
	case AssertionApplyError:
		return "APPLY_ERR"
	default:
		return "?"
	}
}

func applyAssertion(machine Machine, assertion *config.Assertion, e *Executor, printPrefix string) (*AssertionResult, error) {
	result := &AssertionResult{Result: AssertionError}
	var err error

	switch assertion.Kind {
	case config.HashMatchAssrt:
		result, err = applyHashAssertion(machine, assertion)
	case config.FileExistsAssrt:
		result, err = applyExistsAssertion(machine, assertion)
	case config.FileNotExistsAssrt:
		result, err = applyNotExistsAssertion(machine, assertion)
	case config.HashFileAssrt:
		result, err = applyHashFileMatchAssertion(machine, assertion)
	case config.RegexMatchAssrt:
		result, err = applyRegexContentsAssertion(machine, assertion)
	default:
		err = errors.New("unknown assertion kind: " + assertion.Kind)
	}

	if err == nil && result.Result == AssertionApplied { //apply the actions
		for _, action := range assertion.Actions {
			err = doAction(machine, assertion, action, e, printPrefix)
			if err == ErrAssertionsFailed {
				result.Result = AssertionFailed
			}
			if err != nil {
				result.Result = AssertionApplyError
				return result, err
			}
		}
	}

	return result, err
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

func applyNotExistsAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	f, err := machine.ReadFile(assertion.FilePath)
	if err != nil && os.IsNotExist(err) {
		return &AssertionResult{Result: AssertionNoop}, nil
	}
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	f.Close()
	return &AssertionResult{Result: AssertionApplied}, nil
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

func applyRegexContentsAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	matched, err := machine.Grep(assertion.FilePath, assertion.Regex)
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	if !matched {
		return &AssertionResult{Result: AssertionApplied}, nil
	}
	return &AssertionResult{Result: AssertionNoop}, nil
}

func applyHashFileMatchAssertion(machine Machine, assertion *config.Assertion) (*AssertionResult, error) {
	// first check file exists
	f, err := machine.ReadFile(assertion.FilePath)
	if err != nil && os.IsNotExist(err) {
		return &AssertionResult{Result: AssertionApplied}, nil
	}
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	f.Close()

	hasher := md5.New()
	localFile, err := os.Open(util.PathSanitize(assertion.BasePath))
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	defer localFile.Close()

	_, err = io.Copy(hasher, localFile)
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}
	localHash := hex.EncodeToString(hasher.Sum(nil))

	hash, err := machine.Hash(assertion.FilePath)
	if err != nil {
		return &AssertionResult{Result: AssertionError}, err
	}

	if hex.EncodeToString(hash) != strings.ToLower(localHash) {
		return &AssertionResult{Result: AssertionApplied}, nil
	}
	return &AssertionResult{Result: AssertionNoop}, nil
}
