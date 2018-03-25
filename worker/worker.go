package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"time"
)

const (
	Workerv1            = "p6wv1"
	WorkerStateInit     = "ini"
	WorkerStateIdle     = "idl"
	WorkerStateReady    = "rdy"
	WorkerStateRunning  = "run"
	WorkerStateStopping = "stp"
	WorkerStateDead     = "ded"

	PermitMaxCount = 999
	PermitMaxTime  = time.Duration(time.Second * 30)
)

type Worker interface {
	String() string
	Concurrency(uint16)
	Prepare(*configuration.Scenario)
	State() <-chan string
	Permit(Permit)
	Stop()
	Destroy()
}

type Permit struct {
	Time  time.Duration
	Count uint64
}
