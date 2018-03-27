package configuration

import (
	// "fmt"
	// "time"
)
/*
func (soc StaticOrchestrationConfig) Description() string {
	if soc.Count == uint64(0) {
		if soc.Time == time.Duration(0) {
			return "Invalid static orchestration"
		}

		return fmt.Sprintf("Send unlimited requests over %s", soc.Time)
	}

	if soc.Time == time.Duration(0) {
		return fmt.Sprintf("Send %s requests", soc.Count)
	}

	return fmt.Sprintf("Send %s requests over %s", soc.Count, soc.Time)
}

func (doc DynamicOrchestrationConfig) Description() string {
	if doc.ErrorRate <= float32(0.0) {
		if doc.ResponseTime == time.Duration(0) {
			return "Invalid dynamic orchestration"
		}

		return fmt.Sprintf("Send requests until the response time is %s", doc.ResponseTime)
	}

	return fmt.Sprintf("Send requests until the error rate is %s or higher", doc.ErrorRate)
}
*/