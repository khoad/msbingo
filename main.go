package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	bin, err := ioutil.ReadFile("examples/1.bin")
	printErrorAndExit(err)

	// Naive check
	if (bin[0] != 0x56 || bin[1] != 0x02) {
		println("encoded file looks invalid, exiting...")
		os.Exit(1)
	}

	println(dictionaryString[0x04])
}

func printErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}