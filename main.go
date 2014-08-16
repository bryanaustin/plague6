
package main

import (
	"os/signal"
	"os"
	"time"
)

type AllKnowingOne struct {
	AppConfig *AppConfig
	Walks []Walk
	Data []*WalkData
}

type WalkContext struct {
	// Status
	Finished uint64                   // Requests finished
	Locusts uint64                    // Requests made
	Swarms int                        // Active threads
	Walk int                          // Walk number

	// Channels
	Adjust chan int                   // Checks response times and adjusts threads
	Admittance chan int               // Start a single locast.
	Aftermath chan bool               // For stopping devistator
	Devastator chan *ResultProcessing // Does maths on results
	Results chan *Result              // Channel for a locast to return a result
	StatInterval <-chan time.Time     // Timeout for printing to console.
	StatSend chan *StatusData         // Channel for sending status updates
	StatEnd chan *StatusData          // Stop the status thread
	Timeout <-chan time.Time          // Timeout to stop sending locasts.
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

	for w := range ako.Walks {
		walk := InitWalk(w)
		walk.Setup()
		if abort := walk.Run(interchan); abort { break }
		if abort := walk.Finish(interchan); abort { break }
		walk.Cleanup(interchan)
	}
}

func InitWalk(w int) *WalkContext {
	wc := new(WalkContext)
	wc.Walk = w
	wc.Admittance = make(chan int)
	wc.Aftermath = make(chan bool)
	wc.Devastator = make(chan *ResultProcessing)
	wc.Results = make(chan *Result)
	wc.StatInterval = time.After(1)
	wc.StatSend = make(chan *StatusData, 1)
	wc.StatEnd = make(chan *StatusData)
	return wc
}

func (wc *WalkContext) Setup() {
	if ako.AppConfig.Slow > 0 {
		wc.Swarms = 1
		// go SwarmMeasure(wc.Walk, wc.Adjust)
		go Swarm(wc.Admittance, wc.Results)
	} else {
		wc.Swarms = ako.AppConfig.Concurrent
		for x := 0; x < ako.AppConfig.Concurrent; x++ {
			go Swarm(wc.Admittance, wc.Results)
		}
		if ako.AppConfig.Time > 0 {
			wc.Timeout = time.After(time.Duration(ako.AppConfig.Time)*time.Second)
		}
	}

	go Devastation(wc.Devastator, wc.Aftermath)
	go StatusPrint(wc.StatSend, wc.StatEnd)
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
			case wc.Admittance <- wc.Walk:
				wc.Locusts++
				if ako.AppConfig.Requests > 0 {
					keepwalking = wc.Locusts < ako.AppConfig.Requests
				}

			// Finished
			case result := <-wc.Results:
				go func() {
					wc.Devastator <- &ResultProcessing{ wc.Walk, result }
				}()
				wc.Finished++

			// Change concurrency
			case amount := <-wc.Adjust:
				if amount < 0 {
					keepwalking = false
				} else {
					for s := 0; s < amount; s++ {
						go Swarm(wc.Admittance, wc.Results)
					}
				}
				wc.Swarms += amount

			// Print status
			case <-wc.StatInterval:
				wc.StatSend <- &StatusData{
					wc.Walk,
					wc.Swarms,
					ako.Data[wc.Walk].Successful, 
					ako.Data[wc.Walk].Failed, 
					ako.Data[wc.Walk].TotalDuration }
				wc.StatInterval = time.After(time.Duration(ako.AppConfig.StatusInterval)*time.Millisecond)
		}
	}
	return false
}

func (wc *WalkContext) Finish(interchan chan os.Signal) bool {
	if wc.Finished < wc.Locusts {
		for finishwalk := true; finishwalk; {
			select {

				// Second interrupt
				case <- interchan:
					Message("Forced shutdown.")
					finishwalk = false
					return true

				// Finished
				case result := <-wc.Results:
					wc.Devastator <- &ResultProcessing{ wc.Walk, result }
					wc.Finished++
					finishwalk = wc.Finished < wc.Locusts
			}
		}
	}
	wc.Aftermath <- true; <-wc.Aftermath
	wc.StatEnd <- &StatusData{
		wc.Walk,
		wc.Swarms,
		ako.Data[wc.Walk].Successful, 
		ako.Data[wc.Walk].Failed, 
		ako.Data[wc.Walk].TotalDuration }
	wc.StatEnd <- nil; <-wc.StatEnd
	return false
}

func (wc *WalkContext) Cleanup(interchan chan os.Signal) {
	for killed := 0; killed < ako.AppConfig.Concurrent; {
		select {
			case <- interchan:
				Message("Aborting cleanup.")
				return
			case wc.Admittance <- -1:
				killed++
		}
	}
}