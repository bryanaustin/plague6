
package main

import (
	"time"
)

type RollingAvarage struct {
	Data  []time.Duration  // Stored values
	index int              // Last location inserted
	used  int              // Number of filled datum
}

func InitRollingAvarage(count int) *RollingAvarage {
	ra := new(RollingAvarage)
	ra.Data = make([]time.Duration, count)
	return ra
}

func (ra *RollingAvarage) Add(datum time.Duration) {
	if ra.used < len(ra.Data) { ra.used++ }
	ra.Data[ra.index] = datum
	ra.index = (ra.index + 1) % len(ra.Data)
}

func (ra *RollingAvarage) Avarage() time.Duration {
	var total time.Duration
	for i := 0; i < ra.used; i++ {
		total += ra.Data[i]
	}
	if total == 0 { return time.Duration(0) }
	return total / time.Duration(ra.used)
}