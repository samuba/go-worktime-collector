package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Todo: linux/mac support

func main() {
	daysToShow := getDaysToShow()
	if daysToShow == 0 {
		return
	}
	fmt.Print("searching eventlog for last ", daysToShow, " days... ")

	log := readSystemEventlog()
	fmt.Print("found ", len(log), " entries \n\n")

	times := extractAllTimes(log)

	for i := 0; i < daysToShow; i++ {
		day := time.Now().AddDate(0, 0, -i)
		dayTimes := filterForDate(times, day)
		if len(dayTimes) == 0 {
			continue
		}

		first := dayTimes[0]
		last := dayTimes[len(dayTimes)-1]
		if i == 0 {
			last = time.Now()
		}

		fmt.Println(first.Format("02.01.2006 (Monday)"))
		if i == 0 {
			fmt.Println(first.Format("15:04"), "—", "now")
		} else {
			fmt.Println(first.Format("15:04"), "—", last.Format("15:04"))
		}
		fmt.Printf("total: %.2fh\n\n", last.Sub(first).Hours())
	}
}

func getDaysToShow() int {
	if len(os.Args) == 1 {
		return 7 // one week should be sensible default
	}
	arg, err := strconv.Atoi(os.Args[1])
	if err != nil || arg <= 0 {
		printHelp()
		return 0
	}
	return arg
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
		times = append(times, time.Local())
	}
	return times
}

func printHelp() {
	fmt.Println("Will output the startup and shutdown time for each day in the past and calculate the uptime.")
	fmt.Println("WARNING: Assumes that the computer is always shutdown before 00:00!\n")
	fmt.Println("Usage: worktime-collector.exe 10")
	fmt.Println("Will display the last 10 days including today.")
}
