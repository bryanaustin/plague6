package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"time"
)

type Local struct {
	// spawners []Spawner
	scenario chan configuration.Scenario
	concurrency chan uint16
	ready chan struct{}
	state string
}

func NewLocal() (l *Local) {
	l = new(Local)
	l.scenario = make(chan configuration.Scenario)
	l.concurrency = make(chan uint16)
	l.ready = make(chan struct{})
	l.state = WorkerStateInit
	go l.main()
	return
}

func (l *Local) main() {
	// Loop though scenarios
	for  {
		var started bool
		l.state = WorkerStateIdle
		for {
			select {
				//<-scenarioChan
				//<-concurrencyChan
				//<-stateReqChan
			}
		}
		l.state = WorkerStateReady
		l.ready <- struct{}

		// Loop though permits
		for {
			if started {
				l.state = WorkerStateRunning
			} else {
				l.state = WorkerStateReady
			}
			select {
				//<-stopChan
				//<-permitChan + started = true
				//<-doneChan + started = false
				//<-concurrencyChan
				//<-stateReqChan
			}
		}
	}
	close(l.concurrency)
	close(l.scenario)
}

func (l Local) String() string {
	return "<local worker>"
}

// Prepare is an thread safe method to prepare this worker for a scenario
func (l *Local) Prepare(s configuration.Scenario) {
	l.scenario <- s
}

func (l *Local) Ready() (<-chan struct{}) {
	return l.ready
}

func (l *Local) Concurrency(c uint16) error {
	return nil
}

func (l *Local) Permit(n uint64, d time.Duration) error {
	return nil
}

func (l *Local) Stop() error {
	return nil
}

func (l *Local) Destroy() error {
	return nil
}