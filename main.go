package main

import (

	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	log := readSystemEventlog()
	times := extractAllTimes(log)

	d := filterForDate(times, time.Now())

	fmt.Println(d[0])
	fmt.Println("start time: ", d[0])
	fmt.Println("end time:   ", d[len(d)-1])
	fmt.Printf("total: %vh \n", 3)

	fmt.Println("lines: ", len(log))
	fmt.Println("times: ", len(times))
}

func filterForDate(times []time.Time, date time.Time) (result []time.Time) {
	for _, tim := range times {
		if sameDay(tim, date) {
			result = append(result, tim)
		}
	}
	return result
}

func sameDay(time1 time.Time, time2 time.Time) bool {
	y1, m1, d1 := time1.Date()
	y2, m2, d2 := time2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func readSystemEventlog() []string {
	fmt.Println("parsing eventlogs...")
	cmdArgs := strings.Fields("wevtutil qe System ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)

	bytea, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error reading Eventlog:", err)
		return []string{}
	}
	return strings.Split(string(bytea[:len(bytea)-1]), "\n")
}

func extractAllTimes(eventlogLines []string) (times []time.Time) {
	for _, line := range eventlogLines {
		timepattern := "TimeCreated SystemTime='"
		start := strings.Index(line, timepattern)
		if start <= 0 {
			continue
		}
		timeStr := line[start+len(timepattern) : start+len(timepattern)+30]
		timeStr = timeStr[:23] + "Z"
		time, _ := time.Parse("2006-01-02T15:04:05.000Z", timeStr)
		times = append(times, time)
	}
	return times
}
