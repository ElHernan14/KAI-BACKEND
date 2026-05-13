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
