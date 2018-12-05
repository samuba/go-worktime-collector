package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("parsing eventlogs...")

	cmdArgs := strings.Fields("wevtutil qe System /c:10")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	go print(stdout)
	cmd.Wait()
}

// to print the processed information when stdout gets a new line
func print(stdout io.ReadCloser) {
	bytea, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("error: ", err)
	}
	line := string(bytea[:])
	timepattern := "TimeCreated SystemTime='"
	start := strings.Index(line, timepattern)
	fmt.Println("line: ", line[start+len(timepattern):start+len(timepattern)+30])
}
