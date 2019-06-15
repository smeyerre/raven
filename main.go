package main

import {
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
}

func main() {
	msgFile := os.Args[1]

	jsonFile, err := os.Open(msgFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened %v", msgFile)
	defer jsonFile.Close()

	// read opened file as byte array
	byteArrary, _ := ioutil.ReadAll(jsonFile)

	var messageFile MessageFile
	json.Unmarshal(byteArray, &messageFile)
}
