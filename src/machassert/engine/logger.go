package engine

import (
	"fmt"
	"machassert/config"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

// Logger is how the status of assertions and runs are communicated.
type Logger interface {
	LogMachineStatus(string, bool, *config.Machine, error)
	LogAssertionStatus(string, string, *config.Assertion, *AssertionResult, error)
	// AuthenticationPrompt is called by the machine if a password is required and auth.Kind = prompt
	AuthenticationPrompt(prompt string) (string, error)
	// KeyboardInteractiveAuth is called by the machine if prompts recieved and auth.Kind = prompt
	KeyboardInteractiveAuth(user, instruction string, questions []string, echos []bool) ([]string, error)
}

// ConsoleLogger implementes the Logger interface by pretty-printing to the terminal.
type ConsoleLogger struct {
	machines       map[string]machineStatus
	assertionInfo  []*assertionInfo
	currentMachine string
	linesPrinted   int

	haveDoneInteractivePrompt bool
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

// KeyboardInteractiveAuth is called by a machine object if authKind = 'prompt', and a keyboard interactive authentication session is initiated by the server.
func (l *ConsoleLogger) KeyboardInteractiveAuth(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	if instruction != "" {
		l.printf("\n%s (%s)", instruction, user)
	}
	for i := range questions {
		if i == 0 && !l.haveDoneInteractivePrompt {
			l.printf("\n")
			l.haveDoneInteractivePrompt = true
		}
		l.printf("\t%s", questions[i])
		r, err := gopass.GetPasswd()
		l.linesPrinted++
		if err != nil {
			return nil, err
		}
		answers = append(answers, string(r))
	}
	return
}

// AuthenticationPrompt is called by a machine object if authKind = 'prompt', and a password is required.
func (l *ConsoleLogger) AuthenticationPrompt(prompt string) (string, error) {
	l.printf("\n%s", prompt)
	pw, err := gopass.GetPasswd()
	l.linesPrinted++
	return string(pw), err
}

// LogMachineStatus is called when a machine's (being asserted against) status changes.
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
				if assertionInfo.err != nil {
					l.printf(" (%v)", assertionInfo.err)
				}
			}
		}

		if l.currentMachine != assertionInfo.machine {
			l.printf(" (%s)", assertionInfo.machine)
		}
		l.printf("\r\n")
	}

}

// LogAssertionStatus is called with assertion information when an assertion changes status.
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
		for i := range l.assertionInfo {
			if l.assertionInfo[i].assertion == assertion {
				l.assertionInfo[i].result = assertionResult
				l.assertionInfo[i].err = err
			}
		}
	}
	l.paint()
}
