package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	USERNAME string = "Sam Meyer-Reed"
)

const (
	HELP          string = "--help"
	SENT_RECEIVED string = "-s"
	WORD_INFO     string = "-w"
	FLOURISH      string = "--flourish"
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

		if os.Args[1] == FLOURISH { // special case when writing flourish file
			flourish(os.Args[2])
			os.Exit(0)
		}
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
	err = json.Unmarshal(byteArray, &messageFile)
	if err != nil {
		fmt.Println(err)
	}

	switch argument {
	case HELP:
		usage()
	case SENT_RECEIVED:
		sentReceived(messageFile)
	case WORD_INFO:
		fmt.Println(messageFile)
		wordInfo(messageFile)
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
		-w
			prints info related to word counts. I.e. average word count sent and received.
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

	fmt.Println("Total messages:", sent+received)
	fmt.Println("	sent:", sent)
	fmt.Println("	receeived:", received)
}

// NOTE: may be inaccuracies due to content being weird with photos are shared links
func wordInfo(file MessageFile) {
	var sentMessages, receivedMessages, sentWords, receivedWords int
	for _, msg := range file.Messages {
		if msg.SenderName == USERNAME {
			if msg.MessageType == GENERIC && len(msg.Photos) == 0 {
				sentMessages++
				fmt.Println(msg.Content)
				fmt.Println(len(strings.Fields(msg.Content)))
				sentWords += len(strings.Fields(msg.Content))
			}
		} else {
			if msg.MessageType == GENERIC && len(msg.Photos) == 0 {
				receivedMessages++
				receivedWords += len(strings.Fields(msg.Content))
			}
		}

		fmt.Println(sentMessages, receivedMessages, sentWords, receivedWords)
	}

	combinedAvg := (sentWords + receivedWords) / (sentMessages + receivedMessages)
	sentAvg := sentWords / sentMessages
	receivedAvg := receivedWords / receivedMessages

	fmt.Println("Total word count:", sentWords+receivedWords)
	fmt.Println("	sent:", sentWords)
	fmt.Println("	received:", receivedWords)
	fmt.Println("\nCombined average words per message:", combinedAvg)
	fmt.Println("	average words per sent message:", sentAvg)
	fmt.Println("	average words per received message:", receivedAvg)
}
