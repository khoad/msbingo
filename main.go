package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	encodedBytes, err := ioutil.ReadFile("examples/1.bin")
	printErrorAndExit(err)

	println(string(encodedBytes))
}

func printErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}