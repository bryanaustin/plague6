package worker

import (
	"github.com/bryanaustin/plague6/configuration"
)

type Local struct {
	// spawners []Spawner
	scenario chan configuration.Scenario
	concurrency chan uint16
	stateIn chan string
	stateOut chan string
}

func NewLocal() (l *Local) {
	l = new(Local)
	l.scenario = make(chan configuration.Scenario)
	l.concurrency = make(chan uint16)
	l.stateIn = make(chan string, 4)
	l.stateOut = make(chan string)
	go l.stateWatch()
	go l.main()
	return
}

func (l *Local) main() {
	// Loop though scenarios
	var concurrent uint16
	for  {
		var started bool
		l.stateIn <- WorkerStateIdle
		for {
			select {
				//<-scenarioChan
				case c := <-l.concurrency:
					l.spawner.ChangeConcurrecy(c)
				//<-stateReqChan
			}
		}
		l.stateIn <- WorkerStateReady

		// Loop though permits
		for {
			select {
				//<-stopChan
				//<-permitChan + l.stateIn <- WorkerStateRunning
				//<-doneChan + l.stateIn <- WorkerStateReady or WorkerStateFinishing
				//<-concurrencyChan
				//<-stateReqChan
			}
		}
	}
}

func (l *Local) stateWatch() {
	for {
		current := WorkerStateInit
		select {
			case ns := <-l.stateIn:
				if ns == "" { return }
				current = ns
			case l.stateOut <- current:
		}
	}
}

func (l Local) String() string {
	return "<local worker>"
}

// Prepare is an thread safe method to prepare this worker for a scenario
func (l *Local) Prepare(s configuration.Scenario) {
	go func(){
		l.scenario <- s
	}()
}

func (l *Local) State() (<-chan string) {
	return l.stateOut
}

func (l *Local) Concurrency(c uint16) {
	go func(){
		l.concurrency <- c
	}()
}

func (l *Local) Permit(p worker.Permit) {

}

func (l *Local) Stop() {
	// ???
}

func (l *Local) Destroy() error {
	close(l.concurrency)
	close(l.scenario)
}