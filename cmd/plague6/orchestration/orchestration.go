package orchestration

type Orchestration interface {
	// ImmutableWorkerAllocation is true when the worker setup is immutable
	ImmutableWorkerAllocation() bool
}