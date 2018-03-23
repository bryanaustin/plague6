package configuration

import (
	"time"
)

const (
	HeaderAdd   = "add"
	HeaderClear = "clr"
	HeaderSet   = "set"
)

type Orchestration interface {
	Description() string
}

type StaticOrchestrationConfig struct {
	Time  time.Duration
	Count uint64
}

type DynamicOrchestrationConfig struct {
	ErrorRate    float32
	ResponseTime time.Duration
}

// Worker is a satellite that can increase load/bandwidth
type WorkerRemote struct {
	Address string
}

type WorkerLocal struct{}

// Scenario one of possibly multiple tests to run
type Scenario struct {
	Description   string
	Concurrency   uint16
	Requests      []Request
	Probabilities []float32
	Orchestration
}

// Request is the actionable part of the Scenario
type Request struct {
	URL           string
	Method        string
	Body          []byte
	HeaderChanges []HeaderChange
}

// HeaderChange is a change that needs made to the default headers Go provides
type HeaderChange struct {
	Type  string
	Key   string
	Value string
}

// Configuration is all information about what this run is supposed to do
type Configuration struct {
	//workerLookup map[string]worker // Is this needed?
	Workers   []interface{}
	Scenarios []Scenario
}

// New will create a new configuration
func New() (c *Configuration) {
	c = new(Configuration)
	c.Workers = make([]interface{}, 0, 1)
	c.Scenarios = make([]Scenario, 0, 1)
	return
}
