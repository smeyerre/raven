package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"time"
	"sync"
	"fmt"
	"path/filepath"
  "sort"
	"syscall"
)

const (
	DAY       int64  = 86400000
)

func parseConversation(config ConfigFile, convoDir string, w *csv.Writer, lock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

  fmt.Println("Parsing conversation:", convoDir)

	msgFilePaths, err := filepath.Glob(convoDir + "/*" + config.MessageFileType)
	if err != nil {
		fmt.Println("Could not find message file. Path:", convoDir)
		return
	}
  sort.Sort(sort.Reverse(sort.StringSlice(msgFilePaths)))

  // Get converstation participant name
  // ================
  var d FlourishMessageFile
  d, err = GetMessageFileFromJson(msgFilePaths[0])
  if err != nil {
    return 
  }

	if len(d.Participants) < 1 {
		fmt.Println("Could not determine participant. Path:", msgFilePaths[0])
		return
	}
	person := d.Participants[0]

	switch len(d.Participants) {
	case 1:
		// do nothing
	case 2:
		if person.Name == config.Username {
			person = d.Participants[1]
		}
	default:
		// TODO: Handle group chats
		return
	}
  // =================

  var fileIndex int = 0
  var messages []FlourishMessage 
  messages = d.Messages

	var count int = 0
	var messageIndex int = len(messages) - 1
	var currentDay int64 = config.StartDate
	var record []string
	record = append(record, person.Name)

	for { // loop through days
		if currentDay > config.EndDate {
			break
		}

		for { // loop through messages
			if messageIndex < 0 { // finished file
        if fileIndex >= len(msgFilePaths) - 1 { // finished all files
          break
        } else { // move to next file
          fileIndex++

          msgFile, err := GetMessageFileFromJson(msgFilePaths[fileIndex])
          if err != nil {
            return 
          }
          messages = msgFile.Messages
          messageIndex = len(messages) - 1
        }
      }

      if messages[messageIndex].Timestamp <= currentDay {
				count++
				messageIndex--
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

func parseSubdirectory(config ConfigFile, subDir string, w *csv.Writer) error {
  fmt.Println("Parsing subdirectory:", subDir)

	var rLimit syscall.Rlimit
	var maxFiles uint64 = config.DefaultMaxFiles
  // get system limit for open files
  err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err == nil {
		maxFiles = rLimit.Cur
	}

  // Get all convo directories in sub directory
	files, err := os.ReadDir(subDir)
	if err != nil {
		return err
	}

	var lock sync.Mutex
  var wg sync.WaitGroup
	var semaphore = make(chan int, maxFiles)
  // For each of the conversation directories, parse it
	for _, f := range files {
		wg.Add(1)
		semaphore <- 1
		go parseConversation(config, subDir + f.Name(), w, &lock, &wg)
		<-semaphore
	}

	wg.Wait()
	w.Flush()
  
  return nil
}

func Flourish(rootDir string, config ConfigFile) error {
	// Create a csv file
	f, err := os.Create("./flourish.csv") // FIXME: filename
	if err != nil {
		return err
	}
	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)

	// write first csv record as labels for columns
	var currentDay int64 = config.StartDate
	var labels []string

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return err
	}

	labels = append(labels, "Name")

	for { // loop through days
		if currentDay > config.EndDate {
			break
		}
		
		t := time.Unix(currentDay/1000, 0).In(loc)
		labels = append(labels, t.Format("January 2006"))
		currentDay += DAY
	}

  // Write the column labels as the first record
	w.Write(labels)

	contents, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

  // Run through each conversation directory
	for _, dir := range config.ConvoDirectoryNames {
    for _, file := range contents {
      if dir == file.Name() + "/" && file.IsDir() {
        parseSubdirectory(config, rootDir + dir, w)
      }
    }
  }

	return nil
}

func GetMessageFileFromJson(filepath string) (FlourishMessageFile, error) {
  var d FlourishMessageFile

  // Reading data from JSON file
  data, err := os.ReadFile(filepath)
  if err != nil {
    fmt.Println("Could not read file. Error: " + err.Error() + " Path:", filepath)
    return d, err
  }
  // Unmarshal JSON data
  err = json.Unmarshal([]byte(data), &d)
  if err != nil {
    fmt.Println("Could not unmarshal json. Path:", filepath)
    return d, err
  }
  return d, nil
}
