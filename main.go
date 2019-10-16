package main

import (
	"flag"
	"os"
)

func main() {

	boshIP := flag.String("b", "", "Bosh Director IP")
	inputCloudConfig := flag.String("c", "", "Cloud Config json file")
	inputBoshVMsOutput := flag.String("v", "", "Bosh VMS output json file")

	flag.Parse()

	if *inputCloudConfig == "" || *inputBoshVMsOutput == "" || *boshIP == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	compute(*inputCloudConfig, *inputBoshVMsOutput, *boshIP)
}
