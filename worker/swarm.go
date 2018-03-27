package worker

import (
	"github.com/bryanaustin/plague6/configuration"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	LocustErrorTypeNone     = ""
	LocustErrorTypeNetwork  = "net"
	LocustErrorTypeServer   = "srv"
	LocustErrorTypeReadBody = "bod"

	BodyBufferSize = 1024 * 10
)

type FlyResult struct {
	Started, Finished time.Time
	ErrorType         string
	BodySize          uint
	Next              *FlyResult
}

type LocustRequest struct {
	Fly chan<- *configuration.Request
}

type Swarm struct {
	once      sync.Once
	Requester chan<- LocustRequest // Locust seeks permission to fly
	Finisher  chan<- struct{}      // Swarm reports that it is done
	Results   chan<- *FlyResult
}

func (s *Swarm) Start() {
	s.once.Do(func() {
		go s.main()
	})
}

func (s *Swarm) main() {
	bodyBuffer := make([]byte, BodyBufferSize)
	flyPermission := make(chan *configuration.Request)

	for {
		s.Requester <- LocustRequest{Fly: flyPermission}
		r := <-flyPermission
		if r == nil {
			break
		}

		hr := PrepareRequest(r)
		result := &FlyResult{Started: time.Now()}
		resp, err := http.DefaultClient.Do(hr)
		result.Finished = time.Now()
		if err != nil {
			continue
		}

		if resp.StatusCode >= 500 {
			result.ErrorType = LocustErrorTypeServer
		}

		for keepgoing := true; keepgoing; {
			b, rerr := resp.Body.Read(bodyBuffer)
			result.BodySize += uint(b)
			if rerr != nil {
				if rerr != io.EOF {
					result.ErrorType = LocustErrorTypeReadBody
				}
				keepgoing = false
			}
		}
		s.Results <- result
	}
	s.Finisher <- struct{}{}
}
