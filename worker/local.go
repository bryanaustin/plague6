package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"github.com/bryanaustin/plague6/worker/distributor"
	"time"
	// "fmt"
)

type Local struct {
	// spawners []Spawner
	scenario    chan *configuration.Scenario
	concurrency chan uint16
	stateIn     chan string
	stateOut    chan string
	results     chan *FlyResult
	permit      chan *Permit
	warn      chan struct{}
	msg chan interface{}
	done chan *FlyResult
}

func NewLocal() (l *Local) {
	l = new(Local)
	l.scenario = make(chan *configuration.Scenario)
	l.concurrency = make(chan uint16)
	l.stateIn = make(chan string, 4)
	l.stateOut = make(chan string)
	l.results = make(chan *FlyResult)
	l.permit = make(chan *Permit)
	l.warn = make(chan struct{}, 1)
	l.msg = make(chan interface{} , 32)
	l.done = make(chan *FlyResult)
	go l.stateWatch()
	go l.main()
	return
}

func (l *Local) main() {
	// Loop though scenarios
	var concurWant, concurHave, concurAlive uint16
	swarmReq := make(chan LocustRequest)
	swarmFin := make(chan struct{})
	resultFunnel := make(chan *FlyResult)
	statsTick := time.Tick(time.Millisecond * time.Duration(500))

	concurChange := func(c uint16) {
		// l.debug(fmt.Sprintf("concurrency change %d", c))
		concurWant = c
		for ; concurWant > concurHave; concurHave++ {
			s := &Swarm{Requester: swarmReq, Finisher: swarmFin, Results:resultFunnel }
			s.Start()
			concurAlive++
		}
	}

	for {
		var isdone bool
		var permitted bool
		var nmore, nwarn uint64
		var scount, fcount uint64
		var scen *configuration.Scenario
		// var distributor distributor.Distributor
		var dister <-chan interface{}
		var warnTimer, limitTimer <-chan time.Time
		var firstResult *FlyResult
		var lastResult *FlyResult

		concurChange(concurWant)

		collectResult := func(fr *FlyResult) {
			if fr == nil {
				return
			}

			if fr.ErrorType == LocustErrorTypeNone {
				scount++
			} else {
				fcount++
			}

			if firstResult == nil {
				firstResult = fr
			}
			if lastResult == nil {
				lastResult = fr
			} else {
				lastResult.Next = fr
			}
			lastResult = fr
		}

		handleLocust := func(sr LocustRequest){
			if concurHave > concurWant {
				sr.Fly <- nil
				concurHave--
				concurAlive--
				return
			}

			if permitted {

				r := (<-dister).(*configuration.Request)
				sr.Fly <- r

				if nmore > 0 {
					nmore--
					if nmore < 1 {
						isdone = true
					}
				}

				if nmore < nwarn {
					nwarn = 0
					l.warn <- struct{}{}
				}
			}
		}

		l.stateIn <- WorkerStateIdle

		// Wait for setup info
		for keepgoing := true; keepgoing; {
			select {
				case scen = <-l.scenario:
					l.stateIn <- WorkerStateReady
					dister = distQueue(setupDist(scen))
					keepgoing = false

				case c := <-l.concurrency:
					concurChange(c)
			}
		}
		l.stateIn <- WorkerStateReady

		p := <-l.permit
		nmore += p.Count
		nwarn = nmore / 2
		limitTimer = time.After(p.Time)
		warnTimer = time.After(p.Time / 2)
		permitted = true
		l.stateIn <- WorkerStateRunning
		// l.debug("post permit")

		// Loop though permits
		for !isdone {
			select {
			case sr := <-swarmReq:
				handleLocust(sr)

			case fr := <-resultFunnel:
				collectResult(fr)
				// l.toParent(fr)

			case <-warnTimer:
				// Notify parent we are about to run out
				select {
				case l.warn <- struct{}{}:
				default:
				}

			case <-limitTimer:
				permitted = false
				// Stop, cleanup

			case <-statsTick:
				l.toParent(&Stats{ Success:scount, Fail:fcount })

			//<-stopChan
			//<-doneChan + l.stateIn <- WorkerStateReady or WorkerStateFinishing + scen = nil
			case c := <-l.concurrency:
				concurChange(c)
			}
		}

		for concurAlive > 0 {
			select {
			case sr := <-swarmReq:
				sr.Fly <- nil
				concurHave--
				concurAlive--
			case fr := <-resultFunnel:
				collectResult(fr)
				// l.toParent(fr)
			}
		}

		l.stateIn <- WorkerStateStopping
		l.done <- firstResult
		//Cleanup
		// concurWant = concurHave = 0
		// for ; concurAlive > 0; concurAlive-- {
		// 	<-swarmFin
		// }
	}
}

func (l *Local) debug(msg string) {
	l.toParent(&configuration.DebugMessage{ Message: msg })
}

func (l *Local) toParent(x interface{}) {
	select {
	case l.msg <- x:
	default:
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
	current := WorkerStateInit
	for {
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
	l.stateIn <- WorkerStateReady
	go func() {
		l.stateIn <- WorkerStateReady
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

func (l *Local) Permit(p *Permit) {
	go func() {
		l.permit <- p
	}()
}

func (l *Local) Warn() (<-chan struct{}) {
	return l.warn
}

func (l *Local) Done() (<-chan *FlyResult) {
	return l.done
}

func (l *Local) Stop() {
	// ???
}

func (l *Local) Messages() (<-chan interface{}) {
	return l.msg
}




func (l *Local) Destroy() {
	close(l.concurrency)
	close(l.scenario)
}
