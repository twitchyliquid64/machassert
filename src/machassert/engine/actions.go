package engine

import (
	"errors"
	"io"
	"machassert/config"
	"os"
)

func doAction(machine Machine, assertion *config.Assertion, action *config.Action) error {
	switch action.Kind {
	case "":
		return nil
	case config.ActionFail:
		return ErrAssertionsFailed
	case config.ActionCopyFile:
		return copyAction(machine, assertion, action)
	default:
		return errors.New("Unrecognised actions kind: " + action.Kind)
	}
}

func copyAction(machine Machine, assertion *config.Assertion, action *config.Action) error {
	input, err := os.Open(action.SourcePath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := machine.WriteFile(action.DestinationPath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
