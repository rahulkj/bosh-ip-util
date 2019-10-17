package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetRaw(file string) []byte {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return raw
}
