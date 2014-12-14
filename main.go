
package main

import (
	"os/signal"
	"os"
	"time"
)

type AllKnowingOne struct {
	AppConfig *AppConfig
	Walks []Walk
}

type WalkContext struct {
	// Status
	Walk int                          // Walk number
	Finished uint64                   // Requests finished
	Processed uint64                  // Results processed
	Started uint64                    // Requests made

	// Objects
	BlackBox *BlackBox                // Record and report damage
	Swarm *Swarm                      // Manage all the workers
	
	// Channels
	Timeout <-chan time.Time          // Timeout to stop sending locasts.
	StatInterval <-chan time.Time     // Timeout for printing to console.
	Report chan *ProgressReport       // Change something up based on performance
}

var ako *AllKnowingOne

func main() {
	ako = new(AllKnowingOne)
	args := ParseArguments()
	if quit := AttemptListen(); quit { return }
	if quit := DigestArgs(args); quit { return } //TODO: Add Interrupt
	if quit := SetupDataAndTest(); quit { return } //TODO: Add Interrupt
	interchan := make(chan os.Signal, 1)
	signal.Notify(interchan, os.Interrupt)

	for wi, wo := range ako.Walks {
		walk := InitWalk(wi, wo)
		walk.Setup()
		if abort := walk.Run(interchan); abort { break }
		if abort := walk.Finish(interchan); abort { break }
		walk.Cleanup(interchan)
	}
}

func InitWalk(index int, w Walk) *WalkContext {
	wc := new(WalkContext)
	wc.Walk = index
	wc.BlackBox = InitBlackBox(string(index))
	wc.Swarm = InitSwarm(&w)
	wc.StatInterval = time.After(1)
	wc.Report = make(chan *ProgressReport, 1)
	return wc
}

func (wc *WalkContext) Setup() {
	go wc.BlackBox.Start()
	go StatusReporter(wc.Report, wc.BlackBox.Status)
	if ako.AppConfig.Slow > 0 {
		wc.Swarm.SpawnLocust(1)
	} else {
		wc.Swarm.SpawnLocust(ako.AppConfig.Concurrent)
		if ako.AppConfig.Time > 0 {
			wc.Timeout = time.After(time.Duration(ako.AppConfig.Time)*time.Second)
		}
	}
}

func (wc *WalkContext) Run(interchan chan os.Signal) bool {
	for keepwalking := true; keepwalking; {
		select {

			// Signal interrupt 
			case <- interchan:
				Message("\nAttempting clean shutdown.")
				return true
			
			// Timer end
			case <-wc.Timeout:
				keepwalking = false

			// Start next
			case wc.Swarm.Allow <- wc.Walk:
				wc.Started++
				if ako.AppConfig.Requests > 0 {
					keepwalking = wc.Started < ako.AppConfig.Requests
				}

			// Finished
			case result := <-wc.Swarm.Results:
				go func() {
					wc.BlackBox.Intake <- result
					wc.BlackBox.DoneWith <- true
				}()
				wc.Finished++

			// Processed
			case <-wc.BlackBox.DoneWith:
				wc.Processed++

			// Print status
			case <-wc.StatInterval:
				wc.BlackBox.Progress <- true
				strrepr := IntToString(wc.Walk)
				wc.Report <- &ProgressReport{ strrepr, wc.Swarm.Locusts }
				wc.StatInterval = time.After(time.Duration(ako.AppConfig.StatusInterval)*time.Millisecond)
		}
	}
	return false
}

func (wc *WalkContext) Finish(interchan chan os.Signal) bool {
	if wc.Processed < wc.Started {
		for cleaning := true; cleaning; {
			select {

				// Second interrupt
				case <-interchan:
					Message("Forced shutdown.")
					cleaning = false
					return true

				// Finished
				case result := <-wc.Swarm.Results:
					go func() {
						wc.BlackBox.Intake <- result
						wc.BlackBox.DoneWith <- true
					}()
					wc.Finished++

				// Processed
				case <-wc.BlackBox.DoneWith:
					wc.Processed++
					cleaning = wc.Processed < wc.Started
			}
		}
	}
	NewLineYall()
	return false
}

func (wc *WalkContext) Cleanup(interchan chan os.Signal) {
	wc.Report <- nil
	wc.Swarm.Exterminate()
	wc.BlackBox.Stop()
}
