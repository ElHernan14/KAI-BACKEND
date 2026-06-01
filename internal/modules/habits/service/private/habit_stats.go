package private

import (
	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	"time"
)

func BuildUserHabitResponse(habit habitsmodel.UserHabit) habitsdto.UserHabitResponse {
	// records := make([]habitsdto.HabitRecordResponse, 0, len(habit.HabitRecords))
	completedToday := false
	totalXP := 0

	for _, record := range habit.HabitRecords {
		if IsToday(record.Fecha) && record.Completado {
			completedToday = true
		}
		if record.Completado {
			totalXP += record.XPGanada
		}

		// records = append(records, habitsdto.HabitRecordResponse{
		// 	ID:              record.ID.String(),
		// 	Date:            record.Fecha,
		// 	Completed:       record.Completado,
		// 	RegisteredValue: record.ValorRegistrado,
		// 	XPEarned:        record.XPGanada,
		// 	CreatedAt:       record.CreatedAt,
		// })
	}

	return habitsdto.UserHabitResponse{
		ID:             habit.ID.String(),
		HabitCatalogID: habit.HabitCatalogID.String(),
		Name:           habit.HabitCatalog.Name,
		Description:    habit.HabitCatalog.Description,
		Category:       habit.HabitCatalog.Category,
		CareType:       habit.HabitCatalog.CareType,
		Difficulty:     habit.HabitCatalog.Difficulty,
		BaseXP:         habit.HabitCatalog.BaseXP,
		Custom:         habit.Personalized,
		Active:         habit.Active,
		StartDate:      habit.StartDate,
		HabitImage:     habit.HabitCatalog.HabitImage,
		CompletedToday: completedToday,
		TotalXP:        totalXP,
		CurrentStreak:  CalculateCurrentStreak(habit.HabitRecords),
		// Records:        records,
	}
}

func IsToday(value time.Time) bool {
	now := time.Now()
	year, month, day := now.Date()
	valueYear, valueMonth, valueDay := value.Date()

	return year == valueYear && month == valueMonth && day == valueDay
}

func CalculateCurrentStreak(records []habitsmodel.HabitRecord) int {
	completedDates := make(map[string]bool)

	for _, record := range records {
		if !record.Completado {
			continue
		}

		completedDates[DateKey(record.Fecha)] = true
	}

	if len(completedDates) == 0 {
		return 0
	}

	today := DateOnly(time.Now())
	yesterday := today.AddDate(0, 0, -1)

	var currentDate time.Time
	switch {
	case completedDates[DateKey(today)]:
		currentDate = today
	case completedDates[DateKey(yesterday)]:
		currentDate = yesterday
	default:
		return 0
	}

	streak := 0
	for completedDates[DateKey(currentDate)] {
		streak++
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	return streak
}

func DateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func DateKey(value time.Time) string {
	return DateOnly(value).Format("2006-01-02")
}

func CalculateTotalXP(
	records []habitsmodel.HabitRecord,
) int {

	total := 0

	for _, record := range records {
		total += record.XPGanada
	}

	return total
}

func CalculateCompletedRecords(
	records []habitsmodel.HabitRecord,
) int {

	total := 0

	for _, record := range records {

		if record.Completado {
			total++
		}
	}

	return total
}
