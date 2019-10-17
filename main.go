package main

import (
	"fmt"
	"log"
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

	result, err := compute(cloudConfigOutputJSON, boshVMsOutput, boshIP)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] Could perform the operation: %s", err.Error()))
	}

	fmt.Println(result)
}
