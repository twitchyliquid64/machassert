package engine

import (
	"fmt"
	"machassert/config"
	"strings"
)

// Logger is how the status of assertions and runs are communicated.
type Logger interface {
	LogMachineStatus(string, bool, *config.Machine, error)
	LogAssertionStatus(string, string, *config.Assertion, *AssertionResult, error)
}

type ConsoleLogger struct {
	machines       map[string]machineStatus
	assertionInfo  []*assertionInfo
	currentMachine string
}

type machineStatus struct {
	machine     *config.Machine
	err         error
	isConnected bool
}

type assertionInfo struct {
	name      string
	specName  string
	machine   string
	assertion *config.Assertion
	result    *AssertionResult
	err       error
}

func (l *ConsoleLogger) LogMachineStatus(name string, isConnected bool, m *config.Machine, err error) {
	if l.machines == nil {
		l.machines = make(map[string]machineStatus)
	}

	l.machines[name] = machineStatus{
		machine:     m,
		isConnected: isConnected,
		err:         err,
	}
	l.currentMachine = name

	l.paint()
}

func sanitizeName(in string) string {
	return strings.Replace(in, " ", "_", -1)
}

func (l *ConsoleLogger) paint() {

	if l.currentMachine == "" {
		return
	}
	fmt.Print("\033[2J\033[0;0H") //reset screen
	if l.machines[l.currentMachine].isConnected && l.machines[l.currentMachine].err == nil {
		fmt.Printf("Running assertions on %s: %s\n", Cyan(l.currentMachine), Green("CONNECTED"))
	} else if l.machines[l.currentMachine].err == nil {
		fmt.Printf("Connecting to %s:", Cyan(l.currentMachine))
	} else {
		fmt.Printf("Connecting to %s: %s (%s)\n", Cyan(l.currentMachine), Yellow("ERROR"), l.machines[l.currentMachine].err)
		return
	}

	for _, assertionInfo := range l.assertionInfo {
		fmt.Printf("  %s.%s: ", sanitizeName(assertionInfo.specName), sanitizeName(assertionInfo.name))
		if assertionInfo.result == nil {
			fmt.Print(Yellow("RUNNING"))
		} else {
			if assertionInfo.result.Result == AssertionNoop {
				fmt.Print(Green(assertionInfo.result.String()))
			} else if assertionInfo.result.Result == AssertionFailed {
				fmt.Print(Red(assertionInfo.result.String()))
			} else if assertionInfo.result.Result == AssertionApplied {
				fmt.Print(Yellow(assertionInfo.result.String()))
			} else {
				fmt.Print(Red(assertionInfo.result.String()))
			}
		}

		if l.currentMachine != assertionInfo.machine {
			fmt.Printf(" (%s)", assertionInfo.machine)
		}
		fmt.Println()
	}

}

func (l *ConsoleLogger) LogAssertionStatus(specName, assertionName string, assertion *config.Assertion,
	assertionResult *AssertionResult, err error) {
	if assertionResult == nil {
		l.assertionInfo = append(l.assertionInfo, &assertionInfo{
			machine:   l.currentMachine,
			assertion: assertion,
			name:      assertionName,
			specName:  specName,
		})
	} else {
		l.assertionInfo[len(l.assertionInfo)-1].result = assertionResult
		l.assertionInfo[len(l.assertionInfo)-1].err = err
	}
	l.paint()
}
