package main

import (
	"flag"
	"fmt"
	"machassert/config"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var targetsFilePathVar string
var modeVar string

func processFlags() {
	flag.StringVar(&targetsFilePathVar, "targets", "", "Path to targets file")
	flag.Parse()
	modeVar = flag.Arg(0)

	if modeVar == "" {
		fmt.Printf("USAGE: %s <mode> <assertion files>\n", os.Args[0])
		os.Exit(1)
	}

	if targetsFilePathVar != "" { //If it is empty, we use the current machine we are on
		if _, err := os.Stat(targetsFilePathVar); err != nil && os.IsNotExist(err) {
			fmt.Printf("Could not stat targets: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func getTargetSpec() (targets *config.MachineSpec) {
	if targetsFilePathVar == "" { //default: current machine
		targets = config.DefaultTarget()
	} else {
		var err error
		targets, err = config.ParseTargetSpecFile(targetsFilePathVar)
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

	switch modeVar {
	case "print":
		spew.Println(targets)
	}
}
