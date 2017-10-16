
package main

import (
	"testing"
	"time"
)

func TestEmpty(t *testing.T) {
	ra := InitRollingAvarage(128)
	avg := ra.Avarage()
	if avg != time.Duration(0) { t.Errorf("Expected average to be 0 got %q", avg) }
}

func TestSome(t *testing.T) {
	var avg time.Duration
	ra := InitRollingAvarage(4)
	ra.Add(time.Duration(5)) // 5
	avg = ra.Avarage()
	if avg != time.Duration(5) { t.Errorf("Expected average to be 5 got %q", avg) }
	ra.Add(time.Duration(11)) // 5 + 11 = 16
	avg = ra.Avarage()
	if avg != time.Duration(8) { t.Errorf("Expected average to be 8 got %q", avg) }
	ra.Add(time.Duration(44)) // 16 + 44 = 60
	avg = ra.Avarage()
	if avg != time.Duration(20) { t.Errorf("Expected average to be 20 got %q", avg) }
	ra.Add(time.Duration(12)) // 60 + 12 = 72
	avg = ra.Avarage()
	if avg != time.Duration(18) { t.Errorf("Expected average to be 18 got %q", avg) }
	ra.Add(time.Duration(13)) // 72 - 5 + 13 = 80
	avg = ra.Avarage()
	if avg != time.Duration(20) { t.Errorf("Expected average to be 20 got %q", avg) }
	ra.Add(time.Duration(155)) // 80 - 11 + 155 = 224
	avg = ra.Avarage()
	if avg != time.Duration(56) { t.Errorf("Expected average to be 56 got %q", avg) }
}