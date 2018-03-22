package worker

import (
	"archive.bryanaustin.name/plague6/configuration"
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
	// Connect() (string, error)
	// Listen() (<-chan string) // TODO: Not going to use string in the long run
	Prepare(configuration.Scenario) error
	Concurrency(uint16) error
	Permit(uint64, time.Duration) error
	Stop() error
	Destroy() error
}