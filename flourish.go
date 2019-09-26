package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	DAY       int64  = 86400000
	START_DAY int64  = 1280721599000
	END_DAY   int64  = 1569470399000
	USER      string = "Sam Meyer-Reed"
)

func flourish(srcFile string) error {
	// Reading data from JSON file
	data, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	// Unmarshal JSON data
	var d FlourishMessageFile
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		return err
	}

	person := d.Participants[0]
	// if person.Name == USER {
	// 	person = d.Participants[1]
	// }

	var m []FlourishMessage = d.Messages

	// Create a csv file
	f, err := os.Create("./flourish.csv") // FIXME: filename
	if err != nil {
		return err
	}
	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)

	var count int = 0
	var i int = len(m) - 1
	var currentDay int64 = START_DAY
	var labels []string
	var record []string

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return err
	}

	labels = append(labels, "Name")
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

		t := time.Unix(currentDay/1000, 0).In(loc)

		labels = append(labels, t.Format("January 02 2006"))
		record = append(record, strconv.Itoa(count))
		currentDay += DAY
	}

	w.Write(labels)
	w.Write(record)
	w.Flush()

	return nil
}
