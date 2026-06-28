package useractivitysynchronization

import (
	"testing"
	"time"
)

func TestShouldShowReturnMessageAfterMoreThanThreeDays(t *testing.T) {
	now := time.Date(2026, time.June, 26, 10, 0, 0, 0, time.Local)
	threeDaysAgo := now.AddDate(0, 0, -3)
	fourDaysAgo := now.AddDate(0, 0, -4)

	if shouldShowReturnMessage(&threeDaysAgo, now) {
		t.Fatal("no debe mostrar regreso exactamente a los tres dias")
	}
	if !shouldShowReturnMessage(&fourDaysAgo, now) {
		t.Fatal("debe mostrar regreso despues de mas de tres dias")
	}
}

func TestShouldShowReturnMessageRequiresPreviousInteraction(t *testing.T) {
	if shouldShowReturnMessage(nil, time.Now()) {
		t.Fatal("un usuario nuevo no debe recibir un mensaje de regreso")
	}
}
