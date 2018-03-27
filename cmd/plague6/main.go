package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/bryanaustin/plague6/cmd/plague6/orchestration"
	"github.com/bryanaustin/plague6/configuration"
	"github.com/bryanaustin/plague6/worker"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	prepareWaitDuration = time.Minute
)

type WorkerMessage struct {
	Id int
	worker.Worker
	Message interface{}
}

type WorkerStats struct {
	Id int
	worker.Worker
	worker.Stats
}

var l *log.Logger
var cw *configuration.Writer

func init() {
	identifier := "plague6"
	if len(os.Args) > 0 {
		identifier = os.Args[0]
	}
	l = log.New(os.Stderr, identifier+": ", log.LstdFlags)
}

func main() {
	checkInputs()
	config, wp := getConfig()
	orch := perpareOrchestrations(config.Scenarios)
	for i := range orch {
		config.Scenarios[i].Orchestration.Description = orch[i].Describe()
	}
	cw = configuration.NewWriter(os.Stdout)
	passConifg(config)

	// Run through scenarios
	for i, s := range config.Scenarios {
		runScenario(s, orch[i], wp)
	}
}

func checkInputs() {
	// Check for args
	if len(os.Args) > 1 {
		l.Fatal("This program expects no arguments")
	}

	// Check stdin status
	instat, err := os.Stdin.Stat()
	if err != nil {
		l.Fatalf("Can't stat stdin, assumed can't read from it: %s", err)
	}

	// Verify stdin is sending
	if (instat.Mode() & os.ModeCharDevice) != 0 {
		l.Fatal("Plague6 is not getting a configuration from stdin. Exiting.")
	}
}

func getConfig() (*configuration.Configuration, []worker.Worker) {
	// Read input
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		l.Fatalf("Error reading stdin: %s", err)
	}

	// Read input int config
	config := new(configuration.Configuration)
	dec := gob.NewDecoder(bytes.NewBuffer(in))
	if err := dec.Decode(config); err != nil {
		l.Fatalf("Error parsing configuration file: %s", err)
	}
	
	// Make sure there is nothing wrong with the configs
	if err := config.SanityCheck(); err != nil {
		l.Fatalf("Configuration sanity check failed: %s", err)
	}

	// Prepare configs for use and connect to workers
	config.Solidify()
	ws, err := worker.ConnectWorkers(config.Workers)
	if err != nil {
		l.Fatalf("Failed to initialize workers: %s", err)
	}

	return config, ws
}

func perpareOrchestrations(ss []*configuration.Scenario) (nos []orchestration.Orchestration) {
	nos = make([]orchestration.Orchestration, len(ss))
	for i := range ss {
		nos[i] = orchestration.Parse(ss[i].Orchestration)
	}
	return
}

func passConifg(c *configuration.Configuration) {
	if err := cw.WriteObj(c, configuration.MsgTypeConfig); err != nil {
		l.Fatalf("Error converting configuration to output format: %s", err)
	}
}

func runScenario(s *configuration.Scenario, o orchestration.Orchestration, wp []worker.Worker) {
	wpcopy := wp[:]
	prepScenario(s, wpcopy)
	// Post worker states and scenario begin

	// l.Print("Pre concurrency")
	if s.Concurrency < uint16(len(wpcopy)) {
		wpcopy = wpcopy[:s.Concurrency]
		// l.Print("Worker pool: %+v", wpcopy)
		for _, w := range wpcopy {
			// l.Print("Worker: %+v", w)
			w.Concurrency(1)
		}
		// Note about unused scenerios
	} else {
		remaning := s.Concurrency
		each := s.Concurrency / uint16(len(wpcopy))
		for i := 0; i < len(wpcopy) - 1; i ++ {
			remaning -= each
			wpcopy[i].Concurrency(each)
		}
		wpcopy[len(wpcopy) - 1].Concurrency(remaning)
	}

	wrkwarn := make(chan worker.Worker)
	wrkmsg := make(chan WorkerMessage)
	wrkdone := make(chan struct{})
	active := len(wpcopy)
	allo := o.InitalAllocation(active)
	// l.Printf("Initial allocation: %+v", allo)
	for i, w := range wp {
		// l.Printf("Initial allocation worker: %+v", allo[i])
		// l.Printf("O continue: %s", o.Continue())
		go func() {
			for {
				<-w.Warn()
				wrkwarn <- w
			}
		}()
		go func(){
			for {
				m := <-w.Messages()
				wrkmsg <- WorkerMessage{ Id:i, Worker:w, Message:m }
			}
		}()
		go func() {
			fr := <-w.Done()
			for ; fr != nil; fr = fr.Next {
				hit := configuration.Hit{ Id:i, Started:fr.Started, Finished:fr.Finished,
					BodySize:fr.BodySize, ErrorType:fr.ErrorType }
				if err := cw.WriteObj(hit, configuration.MsgTypeHit); err != nil {
					l.Print("Error writing hit: %s", err)
				}
			}
			wrkdone <- struct{}{}
		}()
		w.Permit(allo[i])
	}

	for o.Continue() || active > 0 {
		select {
			case <-wrkdone:
				active--

			case w := <-wrkwarn:
				if p := o.SingleAllocation(); p != nil {
					w.Permit(p)
				}

			case m := <-wrkmsg:
				switch m.Message.(type) {
				case *worker.Stats:
					s := m.Message.(*worker.Stats)
					stats := configuration.WorkerStats{ Id:m.Id, Success:s.Success, Fail:s.Fail }
					if err := cw.WriteObj(stats, configuration.MsgTypeWorkerStats); err != nil {
						l.Print("Error writing stats: %s", err)
					}

				case *configuration.DebugMessage:
					dm := m.Message.(*configuration.DebugMessage)
					dm.Id = m.Id
					if err := cw.WriteObj(dm, configuration.MsgTypeDebugMessage); err != nil {
						l.Print("Error writing debug message: %s", err)
					}

				default:
					l.Printf("Unknown worker message: %+v", m.Message)
				}
		}
	}
}

func prepScenario(s *configuration.Scenario, wp []worker.Worker) {
	readyChan := make(chan error)
	for _, w := range wp {
		w.Prepare(s)
		go waitReady(w, readyChan)
	}

	var errored bool
	for range wp {
		if rr := <-readyChan; rr != nil {
			l.Print("Problem while the workers were preparing: " + rr.Error())
			errored = true
		}
	}

	if errored {
		os.Exit(2)
	}
}

func waitReady(w worker.Worker, rc chan error) {
	stateChan := w.State()
	timeoutChan := time.After(prepareWaitDuration)
	for {
		select {
		case s := <-stateChan:
			if s == worker.WorkerStateReady {
				rc <- nil
				return
			} else {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		case <-timeoutChan:
			rc <- fmt.Errorf("worker %s, took too long to be ready (%s)", w, prepareWaitDuration)
			return
		}
	}
}
