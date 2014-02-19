
package main

import (
	"net/http"
	"fmt"
)

type StepResponse struct {
	Problem error
	Transferred uint64
}

func AttemptListen() bool {
	if len(ako.AppConfig.Listen) > 0 {
		Message("Listening on %q.\n", ako.AppConfig.Listen)
		if lerr := ListenMaster(ako.AppConfig.Listen); lerr != nil {
			Message("Error: %s", lerr)
		}
		return true
	}
	return false
}

func ListenMaster(addr string) error {
	return nil
}

func (hs *HttpStep) Compile() error {
	var err error
	if len(hs.Method) < 1 {
		hs.Method = "GET"
	}
	hs.Request, err = http.NewRequest(hs.Method, hs.Url, nil)
	return err
}

func (hs *HttpStep) Run() *StepResponse {
	var readiter int
	var erriter error
	var totalsize uint64
	client := new(http.Client)
	readbuff := make([]byte, 5120)
	result, reqerr := client.Do(hs.Request)
	if reqerr != nil { return &StepResponse{ reqerr, 0 } }
	for ; erriter != nil; readiter, erriter = result.Body.Read(readbuff) {
		totalsize += uint64(readiter)
	}
	result.Body.Close()
	if result.StatusCode > 399 {
		return &StepResponse{ fmt.Errorf("Got http response: %s", result.Status), totalsize }
	}
	return &StepResponse{ nil, totalsize }
}
