package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("parsing eventlogs...")
	cmdArgs := strings.Fields("wevtutil qe System /c:100")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)

	bytea, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error: ", err)
	}

	lines := strings.Split(string(bytea[:len(bytea)-1]), "\n")

	for _, line := range lines {
		tim, _ := exctractTime(line)
		fmt.Println(tim)
	}

	fmt.Println(len(lines))
}

func exctractTime(str string) (tim time.Time, err error) {
	timepattern := "TimeCreated SystemTime='"
	start := strings.Index(str, timepattern)
	if start <= 0 {
		return time.Now(), errors.New("no time found")
	}
	timeStr := str[start+len(timepattern) : start+len(timepattern)+30]
	timeStr = timeStr[:len(timeStr)-7] + "Z"
	return time.Parse("2006-01-02T15:04:05.000Z", timeStr)
}
