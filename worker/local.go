package worker

import (
	"archive.bryanaustin.name/plague6/configuration"
	"time"
)

type Local struct {
	// spawners []Spawner
	scenario chan configuration.Scenario
	concurrency chan uint16
}

func NewLocal() (l *Local) {
	l = new(Local)
	l.concurrency = make(chan uint16)
	l.scenario = make(chan configuration.Scenario)
	go l.main()
	return
}

func (l *Local) main() {
	// Loop though scenarios
	for  {
		// set state idle
		for {
			select {
				//<-scenarioChan
				//<-concurrencyChan
			}
		}
		// set state ready

		// Loop though permits
		for {
			// set state running
			select {
				//<-stopChan
				//<-permitChan
				//<-doneChan
				//<-concurrencyChan
			}
		}
	}
	close(l.concurrency)
	close(l.scenario)
}

// Prepare is an thread safe method to prepare this worker for scenario
func (l *Local) Prepare(s configuration.Scenario) error {
	l.scenario <- s
}

func (l *Local) Concurrency(c uint16) error {

}

func (l *Local) Permit(n uint64, d time.Duration) error {

}

func (l *Local) Stop() error {

}

func (l *Local) Destroy() error {

}