
package main

import (
	"time"
)

// Survives the plague, tells you what happened
type BlackBox struct {
	Identifier string
	Successful uint64
	Failed uint64
	HighestDuration time.Duration
	LowestDuration time.Duration //TODO: Make inifinty?
	TotalDuration time.Duration
	Started time.Time
	RealDuration time.Duration
	TotalTransferred uint64

	DoneWith chan bool
	Intake chan *Result
	Progress chan bool
	Status chan *BoxReport
}

type BoxReport struct {
	Identifier string
	Successful uint64
	Failed uint64
	Duration time.Duration
}

func InitBlackBox(identifier string) *BlackBox {
	box := new(BlackBox)
	box.Identifier = identifier
	box.DoneWith = make(chan bool)
	box.Intake = make(chan *Result)
	box.Progress = make(chan bool, 1)
	box.Status = make(chan *BoxReport)
	return box
}

// WARNING: Since there is only have a single routine reading
// and processing results, it could fall behind and then adaptive
// code will be meaningless.
// It is assumed the test will always take longer than crunching
// the results.
func (bb *BlackBox) Start() {
	bb.Started = time.Now()
	for {
		select {
			case <-bb.Progress:
				bb.SendReport()
			default:
		}

		select{
			case <-bb.Progress:
				bb.SendReport()
			case result := <- bb.Intake:
				if result == nil { return }
				bb.Process(result)
		}
	}
}

func (bb *BlackBox) SendReport() {
	bb.Status <- &BoxReport{
		bb.Identifier,
		bb.Successful,
		bb.Failed,
		time.Now().Sub(bb.Started) }
}

func (bb *BlackBox) Stop() {
	bb.RealDuration = time.Now().Sub(bb.Started)
	bb.Intake <- nil
}

func (bb *BlackBox) Process(result *Result) {
	if result.Problem != nil {
		bb.Failed++
	} else {
		bb.Successful++
	}
	if result.Duration < bb.LowestDuration || bb.LowestDuration == 0 {
		bb.LowestDuration = result.Duration
	}
	if result.Duration > bb.HighestDuration {
		bb.HighestDuration = result.Duration	
	}
	bb.TotalDuration += result.Duration
	bb.TotalTransferred += result.Transferred
}
