package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

const (
	USERNAME string = "Sam Meyer-Reed"
)

const (
	HELP string = "--help"
	SENT_RECEIVED string = "-s"
)

func main() {
	var msgFile string
	var argument string

	switch len(os.Args) {
	case 1:
		fmt.Println("Error: No input file provided.")
		usage()
		os.Exit(1)
	case 2:
		msgFile = os.Args[1]
		argument = SENT_RECEIVED
	case 3:
		msgFile = os.Args[1]
		argument = os.Args[2]
	default:
		fmt.Println("Error: Unanticipated number of arguments.")
		usage()
		os.Exit(1)
	}

	jsonFile, err := os.Open(msgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	// read opened file as byte array
	byteArray, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var messageFile MessageFile
	json.Unmarshal(byteArray, &messageFile)

	switch argument {
	case HELP:
		usage()
	case SENT_RECEIVED:
		sentReceived(messageFile)
	default:
		fmt.Println("Error: unknown argument. Exiting.")
		usage()
		os.Exit(1)
	}

	os.Exit(0)
}

func usage() {
	var usageMessage string = `
	Usage:

	raven input-file [command]


	command list:
		--help
			prints this usage information.
		-s
			returns two integers, the total sent messages, and the total received messages.
	`

	fmt.Println(usageMessage)
}

func sentReceived(file MessageFile) {
	var sent, received int
	for _, msg := range file.Messages {
		if msg.SenderName == USERNAME {
			sent++
		} else {
			received++
		}
	}

	fmt.Println("Total messages:", sent + received)
	fmt.Println("	sent:", sent)
	fmt.Println("	receeived:", received)
}

