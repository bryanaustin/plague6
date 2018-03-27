package main

import (
	"github.com/bryanaustin/plague6/configuration"
	"os"
	"fmt"
)

var config *configuration.Configuration

func main() {
	r := configuration.NewReader(os.Stdin)
	// headerbuffer := [8]byte{}
	for keepgoing := true; keepgoing; {
		x := r.WaitRead()
		switch x.(type) {
			case *configuration.Configuration:
				descConfig(x.(*configuration.Configuration))
			case *configuration.WorkerStats:
				showWorkerStats(x.(*configuration.WorkerStats))
			case *configuration.DebugMessage:
				fmt.Println("Debug message: " + x.(*configuration.DebugMessage).Message)
			case *configuration.Hit:
				showHit(x.(*configuration.Hit))
			case error:
				fmt.Println(fmt.Sprintf("Received error: %s", x))
			case nil:
				keepgoing = false
			default:
				fmt.Println(fmt.Sprintf("Unknown message: %+v", x))
		}
	}
}

func descConfig(x *configuration.Configuration) {
	config = x
	fmt.Println("Configuration:")
	fmt.Println("  Workers:")
	for i, w := range x.Workers {
		switch w.Type {
			case configuration.WorkerTypeLocal:
				fmt.Println(fmt.Sprintf("    %d: Local", i))
			case configuration.WorkerTypeRemote:
				fmt.Println(fmt.Sprintf("    %d: Remote %s", i, w))
			default:
				fmt.Println("    %d: Unknown Worker", i)
		}
	}
	fmt.Println("  Scenarios:")
	for i, s := range x.Scenarios {
		fmt.Println(fmt.Sprintf("    %d: %s", i, s.Description))
		fmt.Println(fmt.Sprintf("      Orchestration: %s", s.Orchestration.Description))
		fmt.Println(fmt.Sprintf("      Concurrency: %d", s.Concurrency))
		fmt.Println("      Requests:")
		for ir, r := range s.Requests {
			per := s.Probabilities[ir] * float32(100.0)
			fmt.Printf("        %d: (%2.1f%%) %s %s", ir, per, r.Method, r.URL)
			if len(r.Body) > 0 {
				fmt.Printf(" (%d byte request body)", len(r.Body))
			}
			fmt.Println()
		}
	}
}

func workerFmt(w *configuration.Worker) string {
	switch w.Type {
		case configuration.WorkerTypeLocal:
			return "<local>"
		case configuration.WorkerTypeRemote:
			return fmt.Sprintf("<remote %s>", w.Address)
	}
	return "<unknown>"
}

func showWorkerStats (ws *configuration.WorkerStats) {
	workerName := "<unknown>"
	if config != nil && ws.Id < len(config.Workers) {
		workerName = workerFmt(config.Workers[ws.Id])
	}
	fmt.Println(fmt.Sprintf("Worker %s: successful %d, failed %d", workerName, ws.Success, ws.Fail))
}

func showHit(h *configuration.Hit) {
	timediff := h.Finished.Sub(h.Started) //TODO: Migrate to the new time/clock code
	fmt.Println(fmt.Sprintf("%s: %d - %s %s", h.Started, h.BodySize, timediff, h.ErrorType))
}