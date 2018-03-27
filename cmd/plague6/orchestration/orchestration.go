package orchestration

import (
	"github.com/bryanaustin/plague6/worker"
	"fmt"
)

type Orchestration interface {
	// ImmutableWorkerAllocation is true when the worker setup is immutable
	ImmutableWorkerAllocation() bool
	InitalAllocation(int) []*worker.Permit
	SingleAllocation() *worker.Permit
	Continue() bool
	Describe() string
}

type CountOrchestration struct {
	Count uint64
}

func (o CountOrchestration) Describe() string {
	return fmt.Sprintf("Make %d requests", o.Count)
}

// Still needed?
func (o CountOrchestration) ImmutableWorkerAllocation() bool {
	return true
}

func (o CountOrchestration) Continue() bool {
	return o.Count > 0 
}

func (o *CountOrchestration) InitalAllocation(n int) (wp []*worker.Permit) {
	wp = make([]*worker.Permit, n)
	thisallo := lowest(worker.PermitMaxCount*uint64(n), o.Count)
	o.Count -= thisallo
	each := thisallo / uint64(n)
	for i := 0; i < len(wp)-1; i++ {
		wp[i] = new(worker.Permit)
		wp[i].Count = each
		wp[i].Time = worker.PermitMaxTime
		thisallo -= each
	}
	wp[n-1] = new(worker.Permit)
	wp[n-1].Count = thisallo
	wp[n-1].Time = worker.PermitMaxTime
	return
}

func (o *CountOrchestration) SingleAllocation() *worker.Permit {
	if o.Count < 1 {
		return nil
	}

	allo := o.Count
	if allo > 16 {
		allo = lowest(o.Count/2, worker.PermitMaxCount)
	}
	o.Count -= allo
	return &worker.Permit{Count: allo, Time: worker.PermitMaxTime}
}

func lowest(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
