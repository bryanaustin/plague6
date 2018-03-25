package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"github.com/bryanaustin/plague6/worker/distributor"
	"time"
)

type Local struct {
	// spawners []Spawner
	scenario    chan *configuration.Scenario
	concurrency chan uint16
	stateIn     chan string
	stateOut    chan string
	results     chan *FlyResult
	permit      chan Permit
}

func NewLocal() (l *Local) {
	l = new(Local)
	l.scenario = make(chan *configuration.Scenario)
	l.concurrency = make(chan uint16)
	l.stateIn = make(chan string, 4)
	l.stateOut = make(chan string)
	l.results = make(chan *FlyResult)
	l.permit = make(chan Permit)
	go l.stateWatch()
	go l.main()
	return
}

func (l *Local) main() {
	// Loop though scenarios
	var concurWant, concurHave, concurAlive uint16
	swarmReq := make(chan LocustRequest)
	swarmFin := make(chan struct{})

	concurChange := func(c uint16) {
		concurWant = c
		for ; concurWant > concurHave; concurHave++ {
			s := &Swarm{Requester: swarmReq, Finisher: swarmFin}
			s.Start()
			concurAlive++
		}
	}

	for {
		var permitted bool
		var nmore uint64
		var scen *configuration.Scenario
		// var distributor distributor.Distributor
		var dister <-chan interface{}
		var warnTimer, limitTimer <-chan time.Time

		l.stateIn <- WorkerStateIdle

		// Wait for setup info
		for scen == nil {
			select {
			case scen = <-l.scenario:
				dister = distQueue(setupDist(scen))

			case c := <-l.concurrency:
				concurChange(c)
			}
		}
		l.stateIn <- WorkerStateReady

		// Loop though permits
		for {
			select {
			case sr := <-swarmReq:
				if concurHave > concurWant {
					sr.Fly <- nil
					concurHave--
				} else {
					if permitted && nmore != 0 {

						r := (<-dister).(*configuration.Request)
						sr.Fly <- r

						if nmore > 0 {
							nmore--
						}
					}
				}

			case p := <-l.permit:
				nmore += p.Count
				limitTimer = time.After(p.Time)
				warnTimer = time.After(p.Time / 2)
				permitted = true

			case <-warnTimer:
				// Notify parent we are about to run out

			case <-limitTimer:
				permitted = false
				// Stop, cleanup

			//<-stopChan
			//<-permitChan + l.stateIn <- WorkerStateRunning
			//<-doneChan + l.stateIn <- WorkerStateReady or WorkerStateFinishing + scen = nil
			case c := <-l.concurrency:
				concurChange(c)
			}
		}

		l.stateIn <- WorkerStateStopping
		//Cleanup
		// concurWant = concurHave = 0
		// for ; concurAlive > 0; concurAlive-- {
		// 	<-swarmFin
		// }
	}
}

func distQueue(d *distributor.Distributor) (dister chan interface{}) {
	dister = make(chan interface{})
	go func() {
		for {
			dister <- d.Pick()
		}
	}()
	return
}

func setupDist(s *configuration.Scenario) (d *distributor.Distributor) {
	d = new(distributor.Distributor)
	d.Options = make([]distributor.Option, len(s.Requests))
	for i, r := range s.Requests {
		d.Options[i] = distributor.Option{Item: r, Target: s.Probabilities[i]}
	}
	return
}

func (l *Local) stateWatch() {
	for {
		current := WorkerStateInit
		select {
		case ns := <-l.stateIn:
			if ns == "" {
				return
			}
			current = ns
		case l.stateOut <- current:
		}
	}
}

func (l Local) String() string {
	return "<local worker>"
}

// Prepare is an thread safe method to prepare this worker for a scenario
func (l *Local) Prepare(s *configuration.Scenario) {
	go func() {
		l.scenario <- s
	}()
}

func (l *Local) State() <-chan string {
	return l.stateOut
}

func (l *Local) Concurrency(c uint16) {
	go func() {
		l.concurrency <- c
	}()
}

func (l *Local) Permit(p Permit) {
	go func() {
		l.permit <- p
	}()
}

func (l *Local) Stop() {
	// ???
}

func (l *Local) Destroy() {
	close(l.concurrency)
	close(l.scenario)
}
