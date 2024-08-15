package main

import (
	"fmt"
	"os"
)

func GetRaw(file string) []byte {

	raw, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return raw
}
