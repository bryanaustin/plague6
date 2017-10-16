package criticalmass

import (
	"testing"
	"time"
)

func TestHistoryDouble(t *testing.T) {
	startTime := time.Time{}
	h := NewHistory(time.Second * 5)
	historyLenCheck(t, h, 0, "Expected history to start with %d entries, got %d")
	historyCapCheck(t, h, 8, "Expected history to start with a capacity of %d entries, got %d")
	
	clock10ms := startTime.Add(time.Millisecond * 10)
	h.Add(Entry{ When:startTime, Dur:time.Millisecond }, clock10ms) // A 1ms req sent at 0, recorded at 10ms
	historyLenCheck(t, h, 1, "Expected history to have %d entries, got %d")
	historyCapCheck(t, h, 8, "Expected history to still have a capacity of %d entries, got %d")

	clock100ms := startTime.Add(time.Millisecond * 100)
	clock50ms := startTime.Add(time.Millisecond * 50)
	h.Add(Entry{ When:startTime, Dur:time.Millisecond }, clock50ms) // A 1ms req sent at 0, recorded at 50ms
	h.Add(Entry{ When:clock10ms, Dur:time.Millisecond }, clock50ms) // A 1ms req sent at 10, recorded at 50ms
	h.Add(Entry{ When:clock50ms, Dur:time.Millisecond }, clock100ms) // A 1ms req sent at 50, recorded at 100ms
	h.Add(Entry{ When:clock50ms, Dur:time.Millisecond * 3 }, clock100ms) // A 3ms req sent at 50, recorded at 100ms
	historyLenCheck(t, h, 5, "Expected history to have %d entries, got %d")
	historyCapCheck(t, h, 8, "Expected history to still have a capacity of %d entries, got %d")

	clock500ms := startTime.Add(time.Millisecond * 500)
	h.Add(Entry{ When:clock50ms, Dur:time.Millisecond }, clock500ms) // A 1ms req sent at 50, recorded at 500ms
	h.Add(Entry{ When:clock100ms, Dur:time.Millisecond }, clock500ms) // A 1ms req sent at 50, recorded at 500ms
	historyLenCheck(t, h, 7, "Expected history to have %d entries, got %d")
	historyCapCheck(t, h, 8, "Expected history to still have a capacity of %d entries, got %d")

	clock750ms := startTime.Add(time.Millisecond * 750)
	h.Add(Entry{ When:clock500ms, Dur:time.Millisecond }, clock750ms) // A 1ms req sent at 500, recorded at 750ms
	historyLenCheck(t, h, 8, "Expected history to have %d entries, got %d")
	historyCapCheck(t, h, 8, "Expected history to still have a capacity of %d entries, got %d")

	clock900ms := startTime.Add(time.Millisecond * 990)
	h.Add(Entry{ When:clock750ms, Dur:time.Millisecond }, clock900ms) // A 1ms req sent at 750, recorded at 900ms
	historyLenCheck(t, h, 9, "Expected history to have %d entries, got %d")
	historyCapCheck(t, h, 16, "Expected history to still have a capacity of %d entries, got %d")

	historyEntryCheck(t, h, 0, startTime, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 1, startTime, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 2, clock10ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 3, clock50ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 4, clock50ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 5, clock50ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 6, clock100ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 7, clock500ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 8, clock750ms, "Expected index %d to to be %d, got %d")
}

func TestHistoryLifecycle(t *testing.T) {
	startTime := time.Time{}
	h := NewHistory(time.Millisecond * 200)

	clock10ms := startTime.Add(time.Millisecond * 10)
	clock50ms := startTime.Add(time.Millisecond * 50)
	h.Add(Entry{ When:clock10ms, Dur:time.Millisecond }, clock50ms) // A 1ms req sent at 10, recorded at 50ms
	historyLenCheck(t, h, 1, "Expected 1st history to have %d entries, got %d")

	clock100ms := startTime.Add(time.Millisecond * 100)
	clock150ms := startTime.Add(time.Millisecond * 150)
	h.Add(Entry{ When:clock50ms, Dur:time.Millisecond }, clock150ms) // A 1ms req sent at 50, recorded at 150ms
	h.Add(Entry{ When:clock100ms, Dur:time.Millisecond }, clock150ms) // A 1ms req sent at 100, recorded at 150ms
	h.Add(Entry{ When:clock100ms, Dur:time.Millisecond }, clock150ms) // A 1ms req sent at 100, recorded at 150ms
	historyLenCheck(t, h, 4, "Expected 2nd history to have %d entries, got %d")

	clock220ms := startTime.Add(time.Millisecond * 220)
	clock230ms := startTime.Add(time.Millisecond * 230)
	h.Add(Entry{ When:clock220ms, Dur:time.Millisecond }, clock230ms) // A 1ms req sent at 220, recorded at 230ms
	historyLenCheck(t, h, 4, "Expected 3rd history to have %d entries, got %d")
	historyStartIndexCheck(t, h, 1, "Expected history start index to be %d, got %d")
	historyEntryCheck(t, h, 1, clock50ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 2, clock100ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 3, clock100ms, "Expected index %d to to be %d, got %d")
	historyEntryCheck(t, h, 4, clock220ms, "Expected index %d to to be %d, got %d")
	
	clock390ms := startTime.Add(time.Millisecond * 390)
	clock400ms := startTime.Add(time.Millisecond * 400)
	h.Add(Entry{ When:clock390ms, Dur:time.Millisecond }, clock400ms) // A 1ms req sent at 390, recorded at 400ms
	historyLenCheck(t, h, 2, "Expected 4th history to have %d entries, got %d")
	historyStartIndexCheck(t, h, 4, "Expected history start index to be %d, got %d")

	clock410ms := startTime.Add(time.Millisecond * 410)
	clock420ms := startTime.Add(time.Millisecond * 420)
	clock430ms := startTime.Add(time.Millisecond * 430)
	clock440ms := startTime.Add(time.Millisecond * 440)
	clock450ms := startTime.Add(time.Millisecond * 450)
	h.Add(Entry{ When:clock400ms, Dur:time.Millisecond }, clock410ms) // A 1ms req sent at 400, recorded at 410ms
	h.Add(Entry{ When:clock410ms, Dur:time.Millisecond }, clock420ms) // A 1ms req sent at 410, recorded at 420ms
	h.Add(Entry{ When:clock420ms, Dur:time.Millisecond }, clock430ms) // A 1ms req sent at 420, recorded at 430ms
	h.Add(Entry{ When:clock430ms, Dur:time.Millisecond }, clock440ms) // A 1ms req sent at 430, recorded at 440ms
	h.Add(Entry{ When:clock440ms, Dur:time.Millisecond }, clock450ms) // A 1ms req sent at 440, recorded at 450ms
	historyLenCheck(t, h, 6, "Expected 5th history to have %d entries, got %d")
	historyStartIndexCheck(t, h, 0, "Expected history start index to be %d, got %d")
}

func TestAvarage(t *testing.T) {
	startTime := time.Time{}
	h := NewHistory(time.Second * 2)
	clock10ms := startTime.Add(time.Millisecond * 10)
	clock20ms := startTime.Add(time.Millisecond * 20)
	clock30ms := startTime.Add(time.Millisecond * 30)
	clock40ms := startTime.Add(time.Millisecond * 40)
	h.Add(Entry{ When:clock20ms, Dur:time.Millisecond * 50 }, clock10ms) // A 50ms req sent at 10, recorded at 20ms
	h.Add(Entry{ When:clock30ms, Dur:time.Millisecond * 100 }, clock20ms) // A 100ms req sent at 20, recorded at 30ms
	h.Add(Entry{ When:clock40ms, Dur:time.Millisecond * 150 }, clock30ms) // A 150ms req sent at 30, recorded at 40ms

	actual := h.Avarage()
	expected := time.Millisecond * 100
	if actual != expected {
		t.Fatalf("Expected to have an average of %s, got %s", expected, actual)
	}
}

func historyCapCheck(t *testing.T, h *history, expected int, message string) {
	actual := cap(h.entries)
	if actual != expected {
		t.Fatalf(message, expected, actual)
	}
}

func historyLenCheck(t *testing.T, h *history, expected int, message string) {
	actual := h.Len()
	if actual != expected {
		t.Fatalf(message, expected, actual)
	}
}

func historyEntryCheck(t *testing.T, h *history, index int, expected time.Time, message string) {
	if h.entries[index].When != expected {
		t.Logf(message, index, expected, h.entries[index].When)
		t.Fail()
	}
}

func historyStartIndexCheck(t *testing.T, h *history, expected int, message string) {
	if h.start != expected {
		t.Fatalf(message, expected, h.start)
	}
}