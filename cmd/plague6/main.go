package main

import (
	"github.com/bryanaustin/plague6/configuration"
	"github.com/bryanaustin/plague6/cmd/plague6/orchestration"
	"github.com/bryanaustin/plague6/worker"
	"io/ioutil"
	"log"
	"os"
	"bytes"
	"encoding/gob"
)

const (
	prepareWaitDuration = time.Minute
)

var l *log.Logger

func init() {
	identifier := "plague6"
	if len(os.Args) > 0 {
		identifier = os.Args[0]
	}
	l = log.New(os.Stderr, identifier+": ", log.LstdFlags)
}

func main() {
	checkInputs()
	config := getConfig()
	os := perpareOrchestrations(config.Scenarios)
	cw := configuration.NewWriter(os.Stdout)
	passConifg(config, cw)

	// Run through scenarios
	for i, s := range config.Scenarios {
		runScenario(s, os, ws)
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

func getConfig() (*configuration.Configuration) {
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

	return config
}

func perpareOrchestrations(ss []configuration.Scenarios) (nos []*orchestration.Orchestration) {
	nos = make([]*orchestration.Orchestration, len(ss))
	for i := range ss {
		nos[i] = orchestration.Parse(ss[i])
	}
	return
}

func passConifg(c *configuration.Configuration, cw *configuration.Writer) {
	// Encode configuration
	sbconf := new(bytes.Buffer)
	sconf := gob.Encoder(sbconf)
	if err := sconf.Encode(c); err != nil {
		l.Fatalf("Error converting configuration to output format: %s", err)
	}

	// Output configuration
	cw.Write(configuration.MsgTypeConfig, sbconf.Bytes())
}

func runScenario(s Scenario, o *orchestration.Orchestration, wp []worker.Worker) {
	prepScenario(s, wp)
	// Post worker states and scenario begin
	// Allocate concurrency & discard unused workers
}

func prepScenario(s Scenario, wp []worker.Worker) {
	readyChan := make(chan error)
	for _, w := range wp {
		w.Prepare(s)
		go waitReady(w, readyChan)
	}

	var errored bool
	for range wp {
		if rr := <-readyChan; rr != nil {
			l.Error("Problem while the workers were preparing: " + rr.String())
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
				if s == WorkerStateReady{
					rc <- nil
					return
				} else {
					<-time.After(time.Duration(100) * time.Millisecond)
				}
			case <-timeoutChan:
				rc <- fmt.Errorf("worker %s, took too long to be ready (%s)", w, prepareWaitDuration)
				return
		}
	}
}