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

func TestCalculateInactiveDaysIgnoresZeroTime(t *testing.T) {
	service := &Service{}
	zero := time.Time{}

	if got := service.calculateInactiveDays(&zero); got != 0 {
		t.Fatalf("inactive days = %d; want 0", got)
	}
}

func TestCalendarDaysBetweenDoesNotUseDurationLimit(t *testing.T) {
	from := time.Date(1600, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.June, 28, 0, 0, 0, 0, time.UTC)

	if got := calendarDaysBetween(&from, to); got <= 106751 {
		t.Fatalf("calendar days = %d; expected a real calendar difference", got)
	}
}

func TestHasInvalidInactiveDays(t *testing.T) {
	if !hasInvalidInactiveDays(106751) {
		t.Fatal("the saturated duration value must be treated as invalid")
	}
	if hasInvalidInactiveDays(30) {
		t.Fatal("a normal inactivity value must remain valid")
	}
}
