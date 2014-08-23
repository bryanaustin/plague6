
package main

import (
	"time"
)

type Swarm struct {
	Locusts int           // Active threads
	Walk *Walk            // Test instructions
	Allow chan int        // Send out a single locast.
	Results chan *Result  // Channel for a locast to return a result
}

type Result struct {
	Problem error
	Duration time.Duration
	Transferred uint64
}


func InitSwarm(w *Walk) *Swarm {
	swarm := new(Swarm)
	swarm.Walk = w
	swarm.Allow = make(chan int)
	swarm.Results = make(chan *Result)
	return swarm
}

func (s *Swarm) SpawnLocust(count int) {
	for i := 0; i < count; i++ {
		go locust(s)
		s.Locusts++
	}
}

// Perform actions for a single locast attack
func locust(s *Swarm) {
	var problem error
	var trans uint64
	for {
		walk := <-s.Allow
		if walk < 0 { return }
		trans = 0
		problem = nil
		started := time.Now()
		for _, step := range s.Walk.Steps {
			result := step.Run()
			if result.Problem != nil {
				problem = result.Problem
				break
			}
			trans += result.Transferred
		}
		duration := time.Now().Sub(started)
		s.Results <- &Result{ problem, duration, trans }
	}
}

func (s *Swarm) Exterminate() {
	for dead := 0; dead < s.Locusts; dead++ {
		s.Allow <- -1
	}
}
