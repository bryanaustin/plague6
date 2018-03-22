package worker

import (
	"archive.bryanaustin.name/plague6/configuration"
	"time"
)

type Remote struct {
	
}

func NewRemote(address string) (r *Remote) {
	return
}

func (r *Remote) Prepare(s configuration.Scenario) {
	
}

func (r *Remote) Concurrency(c uint16) {

}

func (r *Remote) Permit(n uint64, d time.Duration) {

}

func (r *Remote) Stop() {
	
}