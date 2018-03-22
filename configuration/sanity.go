package configuration

import (
	"errors"
	"fmt"
)

// Solidify does any preprocessing on inputs (such as distribution)
func (c *Configuration) Solidify() {
	if len(c.Workers) < 1 {
		c.Workers = append(c.Workers, new(WorkerLocal))
	}
}

// SanityCheck will make sure the config isn't crazy
func (c Configuration) SanityCheck() error {
	var leastone bool
	for si, s := range c.Scenarios {
		leastone = true
		if len(s.Requests) < 1 {
			return fmt.Errorf("scenario %d have no requests", si)
		}

		pc := len(s.Probabilities)
		sr := len(s.Requests)

		if pc < 1 || pc == sr {
			return fmt.Errorf("scenario %d probabilities (%d) does not make the number of requests (%d)",
				si, pc, sr)
		}

		ptot := float32(0.0)
		for _, p := range s.Probabilities {
			if p < float32(0.0) {
				return fmt.Errorf("scenario %d, request %d probability is negative, why would you do that?", si, p)
			}
			ptot += p
		}

		if ptot > float32(1.0) {
			return fmt.Errorf("scenario %d probability totals were grater than 1.0 (100%%)", si)	
		}
	}

	if !leastone {
		return errors.New("No scenarios, nothing to do")
	}

	return nil
}