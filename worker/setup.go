package worker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bryanaustin/plague6/configuration"
	"io/ioutil"
	"net/http"
)

var (
	UnknownWorkerType = errors.New("unknown worker type")
)

func ConnectWorkers(ws []*configuration.Worker) (nws []Worker, ferr error) {
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

func ConnectWorker(w *configuration.Worker) (nw Worker, err error) {
	switch w.Type {
	case configuration.WorkerTypeLocal:
		nw = NewLocal()
	case configuration.WorkerTypeRemote:
		// rc := w.(configuration.WorkerRemote)
		// nw =  NewRemote(rc.Address)
		err = errors.New("Remote worker not implemented")
	default:
		err = UnknownWorkerType
	}
	return
}

func PrepareRequest(r *configuration.Request) (hr *http.Request) {
	hr = new(http.Request)
	hr.Method = r.Method
	hr.URL = r.ParsedURL
	if r.Body != nil {
		hr.Body = ioutil.NopCloser(bytes.NewBuffer(r.Body))
	}
	if r.HeaderChanges != nil {
		hr.Header = http.Header(make(map[string][]string))
		for _, hc := range r.HeaderChanges {
			switch hc.Type {
			case configuration.HeaderAdd:
				hr.Header.Add(hc.Key, hc.Value)
			case configuration.HeaderClear:
				hr.Header.Del(hc.Key)
			case configuration.HeaderSet:
				hr.Header.Set(hc.Key, hc.Value)
			}
		}
	}
	return
}
