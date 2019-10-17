package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
)

func getDetailsFromBosh() (string, string) {
	cmdName := "bosh"
	cmdArgs := []string{"vms", "--json"}

	boshVMsOutput := executeCommand(cmdName, cmdArgs)

	cmdArgs = []string{"cloud-config"}
	cloudConfigOutput := executeCommand(cmdName, cmdArgs)

	cloudConfigOutputJSON := convertYmlToJSON(cloudConfigOutput)

	return cloudConfigOutputJSON, boshVMsOutput
}

func executeCommand(cmdName string, cmdArgs []string) string {
	var cmdOut []byte
	var err error

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running bosh vms command: ", err)
		os.Exit(1)
	}
	return string(cmdOut)
}

func convertYmlToJSON(cloudConfigOutput string) string {
	j, err := yaml.YAMLToJSON([]byte(cloudConfigOutput))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Conversion error: ", err)
		os.Exit(1)
	}

	return string(j)
}
