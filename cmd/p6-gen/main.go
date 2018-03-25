package main

import (
	"encoding/gob"
	"github.com/bryanaustin/plague6/configuration"
	"log"
	"os"
)

var l *log.Logger

func init() {
	identifier := "p6-gen"
	if len(os.Args) > 0 {
		identifier = os.Args[0]
	}
	l = log.New(os.Stderr, identifier+": ", log.LstdFlags)
}

func main() {
	if len(os.Args) < 2 {
		l.Fatalf("Need a URL and an argument")
	}

	c := configuration.New()
	r := &configuration.Request{ URL:os.Args[1] }
	s := &configuration.Scenario{Description: "This is a desc",
		Concurrency: 1,
		Orchestration: &configuration.StaticOrchestrationConfig{ Count: 10 },
		Requests:[]*configuration.Request{r}}
	c.Scenarios = append(c.Scenarios, s)
	enc := gob.NewEncoder(os.Stdout)
	if err := enc.Encode(c); err != nil {
		l.Fatalf("Error encoding configuration: %s", err)
	}
}
