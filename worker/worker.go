package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"time"
)

const (
	Workerv1 = "p6wv1"
	WorkerStateInit = "ini"
	WorkerStateIdle = "idl"
	WorkerStateReady = "rdy"
	WorkerStateRunning = "run"
	WorkerStateStopping = "stp"
	WorkerStateDead = "ded"
)

type Worker interface {
	String() string
	// Connect() (string, error)
	// Listen() (<-chan string) // TODO: Not going to use string in the long run
	Concurrency(uint16) error
	Prepare(configuration.Scenario) error
	Ready() (<-chan struct{})
	Permit(uint64, time.Duration) error
	Stop() error
	Destroy() error
}