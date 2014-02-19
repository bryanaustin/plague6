
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

var ako *AllKnowingOne

func main() {
	ako = new(AllKnowingOne)
	args := ParseArguments()
	if quit := AttemptListen(); quit { return }
	if quit := DigestArgs(args); quit { return } //TODO: Add Interrupt
	if quit := SetupDataAndTest(); quit { return } //TODO: Add Interrupt
	results := make(chan *Result)
	admittance := make(chan int)
	interchan := make(chan os.Signal, 1)
	signal.Notify(interchan, os.Interrupt, os.Kill)

	for x := 0; x < ako.AppConfig.Concurrent; x++ {
		go Swarm(admittance, results)
	}

	for w := range ako.Walks {
		var abort bool
		var timeout <-chan time.Time
		walkprocessor := make(chan *Result, 64) //This should be more than enough
		if ako.AppConfig.Time > 0 {
			timeout = time.After(time.Duration(ako.AppConfig.Time)*time.Second)
		} else { timeout = nil }
		go Devastation(w, walkprocessor)
		
		for keepwalking := true; keepwalking; {
			select {
				case <- interchan:
					Message("Attempting clean shutdown.")
					keepwalking = false
				case <- timeout:
					keepwalking = false
				case admittance <- w:
					ako.Data[w].Locusts++
					if ako.AppConfig.Requests > 0 {
						keepwalking = ako.Data[w].Locusts < ako.AppConfig.Requests
					}
				case result := <-results:
					walkprocessor <- result
					ako.Data[w].Finished++
				//TODO: Adaptive case
				//TODO: Realtime feedback
			}
		}

		if ako.Data[w].Finished < ako.Data[w].Locusts {
			for finishwalk := true; finishwalk; {
				select {
					case <- interchan:
						Message("Forced shutdown.")
						finishwalk = false
						abort = true
					case result := <-results:
						walkprocessor <- result
						ako.Data[w].Finished++
						finishwalk = ako.Data[w].Finished < ako.Data[w].Locusts
				}
			}
		}

		walkprocessor <- nil
		if abort { return }
	}

	for killed := 0; killed < ako.AppConfig.Concurrent; {
		select {
			case <- interchan:
				Message("Aborting cleanup.")
				return
			case admittance <- -1:
				killed++
		}
	}
	
	PrintResults()
}


func no(thing interface{}) {}