
package main

import (
	"fmt"
	"strconv"
)

type ProgressReport struct {
	Walk string
	Concurrency int
}

func DebugMessage(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func Message(message string, args ... interface{}) {
	if !ako.AppConfig.Quiet {
		fmt.Printf(message + "\n", args...)
	}
}

func StatusReporter(prog <-chan *ProgressReport, box <-chan *BoxReport) {
	for {
		progress := <-prog
		if progress == nil { return }
		PrintStatus(progress, <-box)
	}
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func PrintStatus(pr *ProgressReport, br *BoxReport) {
	totalrequests := br.Successful + br.Failed
	if totalrequests > 0 {
		fmt.Printf("\rWalk %s: Concurrent: %d, Successful: %d, Failed: %d, Average Response: %s         ",
			pr.Walk, pr.Concurrency, br.Successful, br.Failed, br.Average)
	} else {
		fmt.Printf("\rWalk %s: Starting...", pr.Walk)
	}
}

func NewLineYall() {
	fmt.Print("\n")
}
