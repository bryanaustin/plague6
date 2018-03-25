package worker

import (
	"fmt"
	"github.com/bryanaustin/plague6/configuration"
	"time"
)

type Remote struct {
	address string
}

func NewRemote(address string) (r *Remote) {
	r = new(Remote)
	r.address = address
	return
}

func (r Remote) String() string {
	return fmt.Sprintf("<remote worker %s>", r.address)
}

func (r *Remote) Prepare(s configuration.Scenario) {

}

func (r *Remote) Concurrency(c uint16) {

}

func (r *Remote) Permit(n uint64, d time.Duration) {

}

func (r *Remote) Stop() {

}
