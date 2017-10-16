package criticalmass

import (
	"time"
)

type avarageData struct {
	Count int
	Sum time.Duration
}

type Gauge struct {
	Total avarageData
	Short history
	Long history
}

type Limits struct {
	MaxShort time.Duration
	MaxLong time.Duration
}

type Entry struct {
	When time.Time
	Dur time.Duration
}

type history struct{
	Limit time.Duration
	*avarageData
	entries []Entry
	start int
	length int
}

type Trend int

const (
	TrendFlat Trend = iota
	TrendSlower // Bigger numbers
	TrendFaster // Smaller numbers
)

var DefualtLimits = Limits{
	MaxShort:time.Duration(time.Second * 5),
	MaxLong:time.Duration(time.Second * 15),
}

var LongLimits = Limits{
	MaxShort:time.Duration(time.Minute),
	MaxLong:time.Duration(time.Minute * 10),
}

func NewHistory(l time.Duration) *history {
	h := new(history)
	h.Limit = l
	h.entries = make([]Entry, 8)
	h.avarageData = new(avarageData)
	return h
}

func (ad avarageData) Avarage() time.Duration {
	if ad.Count < 1 { return time.Duration(0) }
	return ad.Sum / time.Duration(ad.Count)
}

func (h *history) Avarage() time.Duration {
	return h.avarageData.Avarage()
}

func (h *history) Add(e Entry, current time.Time) {
	h.cleanup(current)
	h.double()
	h.entries[h.appendIndex()] = e
	h.length++
	h.avarageData.Count = h.length
}

func (h *history) double() {
	if cap(h.entries) == h.length {
		oldCap := cap(h.entries)
		newCap := oldCap * 2
		ne := make([]Entry, newCap)
		insertIndex := (oldCap - 1) - h.start
		copy(ne, h.entries[h.start:oldCap])
		copy(ne[insertIndex:], h.entries[:h.start])
		h.entries = ne
		h.start = 0
		h.length = oldCap
	}
}

func (h *history) cleanup(from time.Time) {
	// fmt.Println("")
	for {
		if h.length < 1 { break } // Empty history
		// if from time minus entry time is less then the limit
		// fmt.Println(fmt.Sprintf("%d - %d = %d < %d", from, h.entries[h.start].When, from.Sub(h.entries[h.start].When), h.Limit))
		if from.Sub(h.entries[h.start].When) < h.Limit {
			// fmt.Println("Stop")
			break
		}
		// If cleanup made it this far, remove this entry
		h.avarageData.Sum -= h.entries[h.start].Dur
		h.start = h.normIndex(h.start + 1)
		h.length--
		h.avarageData.Count = h.length
	}
}

func (h *history) appendIndex() int {
	return h.normIndex(h.start + h.length)
}

func (h *history) normIndex(i int) int {
	// math.Mod requires int64 conversion. Fuck that.
	ni := i
	c := cap(h.entries) - 1
	if ni > c { ni -= c }
	return ni
}

func (h *history) Len() int {
	return h.length
}

func (g Gauge) Trend() Trend {
	//TODO: This
	return TrendFlat
}

