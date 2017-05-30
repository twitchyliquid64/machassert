package engine

import (
	"errors"
	"machassert/config"
	"machassert/machine"
)

// Executor stores state/configuration for applying assertions to targets.
type Executor struct {
	machines   *config.MachineSpec
	assertions []*config.AssertionSpec
	logger     Logger
}

// New creates a new executor.
func New(machines *config.MachineSpec, assertions []*config.AssertionSpec) *Executor {
	return &Executor{
		machines:   machines,
		assertions: assertions,
		logger:     &ConsoleLogger{},
	}
}

// Run applies the assertions in the executor to the machines it knows about.
func (e *Executor) Run() error {
	for name, machine := range e.machines.Machine {
		e.logger.LogMachineStatus(name, false, machine, nil)
		m, err := connect(name, machine)
		e.logger.LogMachineStatus(name, true, machine, err)
		if err != nil {
			return err
		}

		for _, assertions := range e.assertions {
			err = e.runAssertionOnMachine(m, assertions)
			if err != nil {
				m.Close()
				return err
			}
		}
		err = m.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) runAssertionOnMachine(machine Machine, assertions *config.AssertionSpec) error {
	for assertionName, assertion := range assertions.Assertions {
		e.logger.LogAssertionStatus(assertions.Name, assertionName, assertion, nil, nil)
		result, err := applyAssertion(machine, assertion)
		e.logger.LogAssertionStatus(assertions.Name, assertionName, assertion, result, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func connect(name string, m *config.Machine) (Machine, error) {
	switch m.Kind {
	case config.KindLocal:
		return machine.ConnectLocal(name, m)
	case config.KindSSH:
		return machine.ConnectRemote(name, m)
	}
	return nil, errors.New("Could not interpret machine kind")
}
