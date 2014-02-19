
package main

import (
	"time"
)

type Result struct {
	Problem error
	Duration time.Duration
	Transferred uint64
}

func Swarm(admin <-chan int, results chan<- *Result) {
	var walk int
	var trans uint64
	var started time.Time
	var problem error
	for {
		walk = <-admin
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

func Devastation(walkindex int, results <-chan *Result) {
	var result *Result
	datafocus := ako.Data[walkindex]
	for {
		result = <-results
		if result == nil { return }
		if result.Problem != nil {
			datafocus.Failed++
		} else {
			datafocus.Successful++
		}
		if result.Duration < datafocus.LowestDuration || datafocus.LowestDuration == 0 {
			datafocus.LowestDuration = result.Duration
		}
		if result.Duration > datafocus.HighestDuration {
			datafocus.HighestDuration = result.Duration	
		}
		datafocus.TotalDuration += result.Duration
		datafocus.TotalTransferred += result.Transferred
	}
}