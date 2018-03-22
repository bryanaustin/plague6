package worker

import (
	"archive.bryanaustin.name/plague6/configuration"
	"errors"
	"fmt"
)

var (
	UnknownWorkerType = errors.New("unknown worker type")
)

func ConnectWorkers(ws []interface{}) (nws []Worker, ferr error) {
	nws = make([]Worker, len(ws))
	for i := range nws {
		nw, err := ConnectWorker(ws[i])
		if err != nil {
			ferr = fmt.Errorf("Could not connect to worker %d (%#v) reason: %s", i, ws[i], err)
			return
		}
		nws[i] = nw
	}
	return
}

func ConnectWorker(w interface{}) (nw Worker, err error) {
	switch w.(type) {
		case configuration.WorkerLocal:
			nw = NewLocal()
		case configuration.WorkerRemote:
			rc := w.(configuration.WorkerRemote)
			nw =  NewRemote(rc.Address)
		default:
			err = UnknownWorkerType
	}
	return
}