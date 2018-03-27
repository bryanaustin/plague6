package orchestration

import (
	"github.com/bryanaustin/plague6/configuration"
)

func Parse(co configuration.Orchestration) Orchestration {
	switch co.Type {
		case configuration.OrchestrationTypeStatic:
			return &CountOrchestration{Count: co.Count}
	}
	// switch co.(type) {
	// case configuration.StaticOrchestrationConfig:
	// 	coc := co.(configuration.StaticOrchestrationConfig)
	// 	return &CountOrchestration{Count: coc.Count}
	// }
	return nil
}
