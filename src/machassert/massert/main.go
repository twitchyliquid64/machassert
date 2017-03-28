package main

import (
	"flag"
	"fmt"
	"machassert/config"
	"os"
)

var targetsFilePathVar string

func processFlags() {
	flag.StringVar(&targetsFilePathVar, "targets", "", "Path to targets file")
	flag.Parse()

	if targetsFilePathVar != "" { //If it is empty, we use the current machine we are on
		if _, err := os.Stat(targetsFilePathVar); err != nil && os.IsNotExist(err) {
			fmt.Printf("Could not stat targets: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func main() {
	processFlags()
	var targets *config.MachineSpec

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
}
