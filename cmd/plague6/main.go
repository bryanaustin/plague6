package main

import (
	"archive.bryanaustin.name/plague6/configuration"
	_ "archive.bryanaustin.name/plague6/cmd/plague6/orchestration"
	"archive.bryanaustin.name/plague6/worker"
	"io/ioutil"
	"log"
	"os"
	"bytes"
	"encoding/gob"
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

	// Encode configuration
	sbconf := new(bytes.Buffer)
	sconf := gob.Encoder(sbconf)
	if err := sconf.Encode(config); err != nil {
		l.Fatalf("Error converting configuration to output format: %s", err)
	}

	// Output configuration
	cw := configuration.NewWriter(os.Stdout)
	cw.Write(configuration.MsgTypeConfig, sbconf.Bytes())

	// Run through scenarios
	for _, s := range config.Scenarios {
		runScenario(s, ws)
	}
}


func runScenario(s Scenario, ws []worker.Worker) {
	var wg sync.WaitGroup
	wg.Add(len(ws))
	for _, w := range ws {
		w.Prepare(s)
		// wg.Done() when ready
	}
	wg.Wait()
}