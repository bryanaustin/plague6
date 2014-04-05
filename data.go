
package main

import (
	"time"
	_ "errors"
)

type WalkData struct {
	Successful uint64
	Failed uint64
	HighestDuration time.Duration
	LowestDuration time.Duration //TODO: Make inifinty?
	TotalDuration time.Duration
	// SegmentTotals [20]time.Duration
	// SegmentStart [20]uint64
	// SegmentFinish [20]uint64
	// SegmentFailed [20]uint64
	Started time.Time
	Duration time.Duration
	TotalTransferred uint64
	//TODO: Add ring for last x request times
}

func (rd *WalkData) AvarageIndividual() time.Duration {
	return rd.TotalDuration / time.Duration(rd.Successful + rd.Failed)
}

// func (rd *WalkData) SegmentAvarage(index int) (time.Duration, error) {
// 	if len(rd.SegmentTotals) < index {
// 		return 0, errors.New("Index outside of range")
// 	}
// 	return rd.SegmentTotals[index] / time.Duration(rd.SegmentFinish[index] + rd.SegmentFailed[index]), nil
// }
