
package main

import (
	"fmt"
	"time"
)

type StatusData struct {
	Successful uint64
	Failure uint64
	TotalDuration time.Duration
}

func Message(message string, args ... interface{}) {
	if !ako.AppConfig.Quiet {
		fmt.Printf(message + "\n", args...)
	}
}

func PrintResults() {
	for d := range ako.Data {
		if len(ako.Data) > 1 {
			Message("Walk %d", d)
		}
		PrintData(ako.Data[d])
	}
}

func PrintData(datas *WalkData) {
	fmt.Printf("Successful: %d\n", datas.Successful)
	fmt.Printf("Failed: %d\n", datas.Failed)
	fmt.Printf("Avarage Response Time: %v\n", datas.AvarageIndividual())
}

func StatusPrint(walk int, statchan <-chan *StatusData, endchan <-chan bool) {
	for {
		select {
			case sd := <-statchan:
				totalrequests := sd.Successful + sd.Failure
				if totalrequests > 0 {
					avarage := sd.TotalDuration / time.Duration(totalrequests)
					fmt.Printf("\rWalk %d: Successful: %d, Failed: %d, Avarage Response: %s",
						walk, sd.Successful, sd.Failure, avarage)
				}
			case <-endchan:
				fmt.Print("\n")
				return
		}
	}
}

func NewLineYall() {
	fmt.Print("\n")
}