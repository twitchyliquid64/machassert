package engine

import (
	"fmt"
	"machassert/config"
	"strconv"
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
	linesPrinted   int
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

func (l *ConsoleLogger) printf(format string, v ...interface{}) {
	out := fmt.Sprintf(format, v...)
	l.linesPrinted += strings.Count(out, "\n")
	fmt.Print(out)
}

func (l *ConsoleLogger) paint() {
	if l.currentMachine == "" {
		return
	}

	// clear lines printed
	fmt.Print(escStart + strconv.Itoa(l.linesPrinted) + "F" + escStart + "0J")
	l.linesPrinted = 0

	if l.machines[l.currentMachine].isConnected && l.machines[l.currentMachine].err == nil {
		l.printf("Running assertions on %s: %s\n", Cyan(l.currentMachine), Green("CONNECTED"))
	} else if l.machines[l.currentMachine].err == nil {
		l.printf("Connecting to %s:", Cyan(l.currentMachine))
	} else {
		l.printf("Connecting to %s: %s (%s)\n", Cyan(l.currentMachine), Yellow("ERROR"), l.machines[l.currentMachine].err)
		return
	}

	for _, assertionInfo := range l.assertionInfo {
		l.printf("  %s.%s: ", sanitizeName(assertionInfo.specName), sanitizeName(assertionInfo.name))
		if assertionInfo.result == nil {
			l.printf(Yellow("RUNNING"))
		} else {
			if assertionInfo.result.Result == AssertionNoop {
				l.printf(Green(assertionInfo.result.String()))
			} else if assertionInfo.result.Result == AssertionFailed {
				l.printf(Red(assertionInfo.result.String()))
			} else if assertionInfo.result.Result == AssertionApplied {
				l.printf(Yellow(assertionInfo.result.String()))
			} else {
				l.printf(Red(assertionInfo.result.String()))
			}
		}

		if l.currentMachine != assertionInfo.machine {
			l.printf(" (%s)", assertionInfo.machine)
		}
		l.printf("\r\n")
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
