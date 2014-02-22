
package main

import (
	"time"
)

type Result struct {
	Problem error
	Duration time.Duration
	Transferred uint64
}

type ResultProcessing struct {
	Walk int
	Result *Result
}

func Swarm(admit <-chan int, results chan<- *Result) {
	var walk int
	var trans uint64
	var started time.Time
	var problem error
	for {
		walk = <-admit
		if walk < 0 { return }
		started = time.Now()
		trans = 0
		problem = nil
		for _, step := range ako.Walks[walk].Steps {
			result := step.Run()
			if result.Problem != nil {
				problem = result.Problem
				break
			}
			trans += result.Transferred
		}
		results <- &Result{ problem, time.Now().Sub(started), trans }
	}
}

func Devastation(results <-chan *ResultProcessing, aftermath chan bool) {
	started := 0
	themaster := make(chan *ResultProcessing)
	ready := make(chan chan *ResultProcessing)
	var result *ResultProcessing

	go TheDevastator(themaster)

	for issue := true; issue; {
		select {
			case result = <-results:
				select {
					case individual := <-ready:
						individual <- result
					default:
						nc := make(chan *ResultProcessing)
						go SomethingDevastating(ready, nc, themaster)
						nc <- result
						started++
				}
			case <-aftermath:
				issue = false
		}
	}
	for killed := 0; killed < started; killed++ {
		readychan := <-ready
		readychan <- nil
	}
	themaster <- nil
	aftermath <- true
}

func SomethingDevastating(ready chan<- chan *ResultProcessing, feed, master chan *ResultProcessing) {
	var process *ResultProcessing
	for {
		process = <-feed
		if process == nil { return }
		master <- process
		ready <- feed
	}
}

func TheDevastator(results <-chan *ResultProcessing) {
	var process *ResultProcessing
	for {
		process = <-results
		if process == nil {
			return
		}
		datafocus := ako.Data[process.Walk]
		if process.Result.Problem != nil {
			datafocus.Failed++
		} else {
			datafocus.Successful++
		}
		if process.Result.Duration < datafocus.LowestDuration || datafocus.LowestDuration == 0 {
			datafocus.LowestDuration = process.Result.Duration
		}
		if process.Result.Duration > datafocus.HighestDuration {
			datafocus.HighestDuration = process.Result.Duration	
		}
		datafocus.TotalDuration += process.Result.Duration
		datafocus.TotalTransferred += process.Result.Transferred
	}
}
