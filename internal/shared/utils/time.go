package utils

import "time"

func Now() time.Time {
	return time.Now().UTC()
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func IsSameCalendarDay(value *time.Time, reference time.Time) bool {
	if value == nil {
		return false
	}

	valueYear, valueMonth, valueDay := value.In(time.Local).Date()
	referenceYear, referenceMonth, referenceDay := reference.In(time.Local).Date()

	return valueYear == referenceYear &&
		valueMonth == referenceMonth &&
		valueDay == referenceDay
}
