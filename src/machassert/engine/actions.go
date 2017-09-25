package engine

import (
	"errors"
	"io"
	"machassert/config"
	"machassert/util"
	"os"
)

func doAction(machine Machine, assertion *config.Assertion, action *config.Action, e *Executor, printPrefix string) error {
	switch action.Kind {
	case "":
		return nil
	case config.ActionFail:
		return ErrAssertionsFailed
	case config.ActionCopyFile:
		return copyAction(machine, assertion, action)
	case config.ActionAssert:
		return assertAction(machine, assertion, action, e, printPrefix)
	default:
		return errors.New("Unrecognised actions kind: " + action.Kind)
	}
}

func assertAction(machine Machine, assertion *config.Assertion, action *config.Action, e *Executor, printPrefix string) error {
	for _, assertionName := range sortAssertions(action.Assertions) {
		assertion := action.Assertions[assertionName]
		e.logger.LogAssertionStatus(printPrefix, assertionName, assertion, nil, nil)
		result, err := applyAssertion(machine, assertion, e, printPrefix+"."+assertionName)
		e.logger.LogAssertionStatus(printPrefix, assertionName, assertion, result, err)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyAction(machine Machine, assertion *config.Assertion, action *config.Action) error {
	output, err := machine.WriteFile(action.DestinationPath)
	if err != nil {
		return err
	}
	defer output.Close()

	input, err := os.Open(util.PathSanitize(action.SourcePath))
	if err != nil {
		return err
	}
	defer input.Close()

	_, err = io.Copy(output, input)
	return err
}
