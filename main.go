package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	Topic     string
	Partition int
	Leader    bool
	Replica   int
	Start     int64
	End       int64
}

var state string
var topic string

func processLine(line string) []*Record {
	if line == "SUMMARY" || line == "CONFIGS" || line == "PARTITIONS" {
		state = line
	}
	if state == "SUMMARY" && strings.HasPrefix(line, "NAME") {
		topic = strings.Fields(line)[1]
	}
	if state == "PARTITIONS" {
		if len(line) > 0 && !(strings.HasPrefix(line, "PA") || strings.HasPrefix(line, "==")) {
			line = strings.Replace(line, "[", "", 1)
			line = strings.Replace(line, "]", "", 1)
			fields := strings.Fields(line)

			replicaCount := len(fields) - 5

			partition, err := strconv.ParseInt(fields[0], 10, 64)
			if err != nil {
				panic(err)
			}

			leader, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				panic(err)
			}

			start, err := strconv.ParseInt(fields[replicaCount+3], 10, 64)
			if err != nil {
				panic(err)
			}

			end, err := strconv.ParseInt(fields[replicaCount+4], 10, 64)
			if err != nil {
				panic(err)
			}

			results := make([]*Record, 0)

			for i := 0; i < replicaCount; i++ {
				replica, err := strconv.ParseInt(fields[3+i], 10, 64)
				if err != nil {
					panic(err)
				}

				results = append(results, &Record{
					Topic:     topic,
					Partition: int(partition),
					Leader:    leader == replica,
					Replica:   int(replica),
					Start:     start,
					End:       end,
				})

			}

			return results
		}
	}

	return nil
}

func processFile(filename string) []*Record {

	results := make([]*Record, 0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		records := processLine(scanner.Text())
		if len(records) > 0 {
			results = append(results, records...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

func main() {
	filename := flag.String("input", "input.txt", "input file to process")
	flag.Parse()

	records := processFile(*filename)
	for _, record := range records {
		fmt.Printf("%v,%v,%v,%v,%v,%v\n",
			record.Topic, record.Partition, record.Leader, record.Replica, record.Start, record.End)
	}
}
