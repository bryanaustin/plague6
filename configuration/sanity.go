package configuration

import (
	"errors"
	"fmt"
	"net/url"
)

// Solidify does any preprocessing on inputs (such as distribution)
func (c *Configuration) Solidify() {
	if len(c.Workers) < 1 {
		c.Workers = append(c.Workers, &Worker{ Type:WorkerTypeLocal })
	}

	for _, s := range c.Scenarios {
		if len(s.Probabilities) < 1 {
			reqcount := len(s.Requests)
			s.Probabilities = make([]float32, reqcount)
			remaining := float32(1.0)
			each := remaining / float32(reqcount)
			for i := 0; i < reqcount-1; i++ {
				s.Probabilities[i] = each
				remaining -= each
			}
			s.Probabilities[reqcount-1] = remaining
		}

		for _, r := range s.Requests {
			if r.Method == "" {
				r.Method = "GET"
			}
		}
	}
}

// SanityCheck will make sure the config isn't crazy
func (c *Configuration) SanityCheck() error {

	if len(c.Scenarios) < 1 {
		return errors.New("No scenarios, nothing to do")
	}

	for si, s := range c.Scenarios {
		if len(s.Requests) < 1 {
			return fmt.Errorf("scenario %d has no requests", si)
		}

		pc := len(s.Probabilities)
		sr := len(s.Requests)

		if pc > 0 && pc != sr {
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

		for ri, r := range s.Requests {
			if r.URL == "" {
				return fmt.Errorf("scenario %d, request %d doesn't have a url", si, ri)
			}

			var err error
			if c.Scenarios[si].Requests[ri].ParsedURL, err = url.Parse(r.URL); err != nil {
				return fmt.Errorf("url (%s) in scenario %d, request %d failed to parse: %s", r.URL, si, ri, err)
			}
			// fmt.Errorf("URL: %s", c.Scenarios[si].Requests[ri].ParsedURL)
		}
	}

	return nil
}
