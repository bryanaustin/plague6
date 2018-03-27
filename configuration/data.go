package configuration

import (
	"net/url"
	"time"
)

const (
	HeaderAdd   = "add"
	HeaderClear = "clr"
	HeaderSet   = "set"

	OrchestrationTypeStatic = ""
	OrchestrationTypeError  = "oer"
	OrchestrationTypeResp   = "rep"

	WorkerTypeLocal = ""
	WorkerTypeRemote = "rmt"
)

type Orchestration struct {
	Description  string
	Type         string
	Time         time.Duration
	Count        uint64
	ErrorRate    float32
	ResponseTime time.Duration
}

type Worker struct {
	Type string
	Address string
}

// Scenario one of possibly multiple tests to run
type Scenario struct {
	Description   string
	Concurrency   uint16
	Requests      []*Request
	Probabilities []float32
	Orchestration
}

// Request is the actionable part of the Scenario
type Request struct {
	URL           string
	ParsedURL     *url.URL
	Method        string
	Body          []byte
	HeaderChanges []HeaderChange
	// Timeout?
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
	Workers   []*Worker
	Scenarios []*Scenario
}

// New will create a new configuration
func New() (c *Configuration) {
	c = new(Configuration)
	c.Workers = make([]*Worker, 0, 1)
	c.Scenarios = make([]*Scenario, 0, 1)
	return
}
