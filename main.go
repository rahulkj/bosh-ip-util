package main

import (
	"fmt"
	"os"
)

func main() {

	envVars := []string{"BOSH_ENVIRONMENT", "BOSH_CLIENT", "BOSH_CLIENT_SECRET", "BOSH_CA_CERT"}

	var missingVars []string
	for i := range envVars {
		_, ok := os.LookupEnv(envVars[i])
		if !ok {
			missingVars = append(missingVars, envVars[i])
		}
	}

	if len(missingVars) > 0 {
		fmt.Println("These environment variables are not set: ", missingVars)
		os.Exit(1)
	}

	boshIP, _ := os.LookupEnv("BOSH_ENVIRONMENT")

	cloudConfigOutputJSON, boshVMsOutput := getDetailsFromBosh()

	compute(cloudConfigOutputJSON, boshVMsOutput, boshIP)
}
