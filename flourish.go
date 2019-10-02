package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"sync"
	"fmt"
	"path/filepath"
	"syscall"
)

const (
	DAY       int64  = 86400000
	START_DAY int64  = 1280721599000
	END_DAY   int64  = 1569470399000
	USER      string = "Sam Meyer-Reed"
	DEFAULT_MAX_FILES uint64 = 1000
)

func parseConversation(convoDir string, w *csv.Writer, lock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	msgFilePaths, err := filepath.Glob(convoDir + "/*.json")
	if err != nil {
		fmt.Println("Could not find message file. Path:", convoDir)
		return
	} else if len(msgFilePaths) != 1 {
		fmt.Printf("Expected 1 possible message file, found %v. Path: %v", len(msgFilePaths), convoDir)
		return
	}

	msgFilePath := msgFilePaths[0]

	// Reading data from JSON file
	data, err := ioutil.ReadFile(msgFilePath)
	if err != nil {
		fmt.Println("Could not read file. Path:", msgFilePath)
		return
	}
	// Unmarshal JSON data
	var d FlourishMessageFile
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		fmt.Println("Could not unmarshal json. Path:", msgFilePath)
		return
	}

	person := d.Participants[0]
	switch len(d.Participants) {
	case 1:
		// do nothing
	case 2:
		if person.Name == USER {
			person = d.Participants[1]
		}
	default:
		// TODO: Handle group chats
		return
	}

	var m []FlourishMessage = d.Messages

	var count int = 0
	var i int = len(m) - 1
	var currentDay int64 = START_DAY
	var record []string
	record = append(record, person.Name)

	for { // loop through days
		if currentDay > END_DAY {
			break
		}

		for { // loop through messages
			if i >= 0 && m[i].Timestamp <= currentDay {
				count++
				i--
			} else {
				break
			}
		}

		record = append(record, strconv.Itoa(count))
		currentDay += DAY
	}

	lock.Lock()
	w.Write(record)
	lock.Unlock()
}

func flourish(rootDir string) error {
	// Create a csv file
	f, err := os.Create("./flourish.csv") // FIXME: filename
	if err != nil {
		return err
	}
	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)

	var currentDay int64 = START_DAY
	var labels []string

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return err
	}

	labels = append(labels, "Name")

	for { // loop through days
		if currentDay > END_DAY {
			break
		}
		
		t := time.Unix(currentDay/1000, 0).In(loc)
		labels = append(labels, t.Format("January 2006"))
		currentDay += DAY
	}

	w.Write(labels)


	var rLimit syscall.Rlimit
	var maxFiles uint64 = DEFAULT_MAX_FILES
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err == nil {
		maxFiles = rLimit.Cur
	}

	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return err
	}

	var lock sync.Mutex
    var wg sync.WaitGroup
	var semaphore = make(chan int, maxFiles)
	for _, f := range files {
		wg.Add(1)
		semaphore <- 1
		go parseConversation(rootDir + f.Name(), w, &lock, &wg)
		<-semaphore
	}

	wg.Wait()
	w.Flush()

	return nil
}
