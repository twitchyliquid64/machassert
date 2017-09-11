package engine

import (
	"errors"
	"machassert/config"
	"machassert/machine"
	"sort"
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

type assertionForSort struct {
	name      string
	assertion *config.Assertion
}

func sortAssertions(assertions map[string]*config.Assertion) []string {
	var out []assertionForSort
	for k, v := range assertions {
		out = append(out, assertionForSort{k, v})
	}
	bo := ByOrder(out)
	sort.Sort(bo)
	return bo.keys()
}

// ByOrder implements sort.Interface for []assertionsForSort based on
// the Order field.
type ByOrder []assertionForSort

func (a ByOrder) Len() int      { return len(a) }
func (a ByOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool {
	if a[i].assertion.Order == a[j].assertion.Order {
		return a[i].name < a[j].name
	}
	return a[i].assertion.Order < a[j].assertion.Order
}
func (a ByOrder) keys() []string {
	out := make([]string, len(a))
	for i := range a {
		out[i] = a[i].name
	}
	return out
}

func (e *Executor) runAssertionOnMachine(machine Machine, assertions *config.AssertionSpec) error {
	for _, assertionName := range sortAssertions(assertions.Assertions) {
		assertion := assertions.Assertions[assertionName]
		e.logger.LogAssertionStatus(assertions.Name, assertionName, assertion, nil, nil)
		result, err := applyAssertion(machine, assertion, e)
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
