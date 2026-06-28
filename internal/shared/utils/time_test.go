package utils

import (
	"testing"
	"time"
)

func TestIsSameCalendarDay(t *testing.T) {
	location := time.FixedZone("test", -3*60*60)
	previousLocation := time.Local
	time.Local = location
	defer func() { time.Local = previousLocation }()

	reference := time.Date(2026, time.June, 28, 22, 0, 0, 0, location)
	sameDay := time.Date(2026, time.June, 28, 1, 0, 0, 0, location)
	previousDay := time.Date(2026, time.June, 27, 23, 59, 0, 0, location)

	if !IsSameCalendarDay(&sameDay, reference) {
		t.Fatal("timestamps del mismo dia deben coincidir")
	}
	if IsSameCalendarDay(&previousDay, reference) {
		t.Fatal("timestamps de dias distintos no deben coincidir")
	}
	if IsSameCalendarDay(nil, reference) {
		t.Fatal("un timestamp nil no debe considerarse sincronizado")
	}
}
