package main

import (
	"fmt"
  "os"
  "io/ioutil"
)

func GetRaw(file string) []byte {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return raw
}
