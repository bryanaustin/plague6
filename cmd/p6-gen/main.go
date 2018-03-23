package main

import (
	"github.com/bryanaustin/plague6/configuration"
	"encoding/gob"
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
	c := configuration.New()
	s := configuration.Scenario{Description: "This is a desc"}
	c.Scenarios = append(c.Scenarios, s)
	enc := gob.NewEncoder(os.Stdout)
	if err := enc.Encode(c); err != nil {
		log.Fatalf("Error encoding configuration: %s", err)
	}
}
