
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
			if d > 0 {
				NewLineYall()
			}
			Message("Walk %d", d)
		}
		PrintData(ako.Data[d])
	}
}

func PrintData(datas *WalkData) {
	fmt.Printf("Successful Steps: %d\n", datas.Successful)
	fmt.Printf("Failed Steps: %d\n", datas.Failed)
	fmt.Printf("Failed Steps: %d\n", datas.Failed)
	fmt.Printf("Avarage Response Time: %v\n", datas.AvarageIndividual())
}

func StatusPrint(walk int, statchan <-chan *StatusData, endchan chan *StatusData) {
	process := func(sd *StatusData) {
		totalrequests := sd.Successful + sd.Failure
		if totalrequests > 0 {
			avarage := sd.TotalDuration / time.Duration(totalrequests)
			fmt.Printf("\rWalk %d: Successful: %d, Failed: %d, Avarage Response: %s",
				walk, sd.Successful, sd.Failure, avarage)
		} else {
			fmt.Printf("\rWalk %d: Starting...", walk)
		}
	}
	var sd *StatusData
	for {
		select {
			case sd = <-statchan:
				process(sd)
			case sd = <-endchan:
				process(sd)
				fmt.Print("\n")
				endchan <- nil
				return
		}
	}
}

func NewLineYall() {
	fmt.Print("\n")
}