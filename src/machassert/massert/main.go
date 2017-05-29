package main

import (
	"errors"
	"flag"
	"fmt"
	"machassert/config"
	"machassert/engine"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var (
	targetsFilePathVar = flag.String("targets", "", "Path to targets file")
	assertionsFiles    []string
	modeVar            string
)

func processFlags() {
	flag.Parse()
	modeVar = flag.Arg(0)

	if modeVar == "" {
		fmt.Printf("USAGE: %s [--targets <target file>] <mode> <assertion files>\n", os.Args[0])
		os.Exit(1)
	}
	assertionsFiles = flag.Args()[1:]

	if *targetsFilePathVar != "" { //If it is empty, we use the current machine we are on
		if _, err := os.Stat(*targetsFilePathVar); err != nil && os.IsNotExist(err) {
			fmt.Printf("Could not stat targets: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func getAssertionsSpecs() ([]*config.AssertionSpec, error) {
	if len(assertionsFiles) == 0 {
		return nil, errors.New("no assertion files provided")
	}
	var out []*config.AssertionSpec
	for _, path := range assertionsFiles {
		assertions, err := config.ParseAssertionsSpecFile(path)
		if err != nil {
			return nil, err
		}
		out = append(out, assertions)
	}
	return out, nil
}

func getTargetSpec() (targets *config.MachineSpec) {
	if *targetsFilePathVar == "" { //default: current machine
		targets = config.DefaultTarget()
	} else {
		var err error
		targets, err = config.ParseTargetSpecFile(*targetsFilePathVar)
		if err != nil {
			fmt.Printf("Err parsing targets: %s\n", err.Error())
			os.Exit(1)
		}
	}
	return targets
}

func main() {
	processFlags()
	targets := getTargetSpec()
	assertions, err := getAssertionsSpecs()
	if err != nil {
		fmt.Println("Err:", err.Error())
		os.Exit(1)
	}

	switch modeVar {
	case "run":
		fallthrough
	case "assert":
		e := engine.New(targets, assertions)
		err = e.Run()
		if err != nil && err == engine.ErrAssertionsFailed {
			fmt.Println(engine.Red("Error") + ": Assertions failed")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println("Err:", err.Error())
			os.Exit(1)
		}

	case "print":
		fmt.Println("Targets:")
		spew.Println(targets)
		fmt.Println()

		fmt.Println("Assertions:")
		spew.Println(assertions)
	}
}
