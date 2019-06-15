package main

import {
	"fmt"
	"os"
}

func main() {
	msgFile := os.Args[1]

	jsonFile, err := os.Open(msgFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened %v", msgFile)
	defer jsonFile.Close()
}
