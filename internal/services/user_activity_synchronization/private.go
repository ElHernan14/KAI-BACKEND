package useractivitysynchronization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) findLastCompletedHabitDate(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*time.Time, error) {

	return s.habitRepository.FindLastCompletedHabitDate(
		ctx,
		tx,
		userID,
	)
}

func (s *Service) calculateInactiveDays(
	lastActivity *time.Time,
) int {
	return calendarDaysBetween(lastActivity, time.Now())
}

func (s *Service) calculateGlobalStreak(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (int, error) {

	dates, err := s.habitRepository.
		FindCompletedHabitDates(
			ctx,
			tx,
			userID,
		)
	if err != nil {
		return 0, err
	}

	if len(dates) == 0 {
		return 0, nil
	}

	streak := 1

	for i := 1; i < len(dates); i++ {

		current := truncateDate(
			dates[i-1],
		)

		next := truncateDate(
			dates[i],
		)

		diff := int(
			current.Sub(next).Hours() / 24,
		)

		if diff == 1 {
			streak++
			continue
		}

		break
	}

	return streak, nil
}

func truncateDate(
	t time.Time,
) time.Time {

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		time.Local,
	)
}

func shouldShowReturnMessage(lastInteraction *time.Time, now time.Time) bool {
	return calendarDaysBetween(lastInteraction, now) >
		returnMessageAfterInactiveDays
}

func calendarDaysBetween(from *time.Time, to time.Time) int {
	if from == nil || from.IsZero() {
		return 0
	}

	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	if !fromDate.Before(toDate) {
		return 0
	}

	return int((toDate.Unix() - fromDate.Unix()) / int64(24*time.Hour/time.Second))
}

func hasInvalidInactiveDays(days int) bool {
	return days < 0 || days > maxSaneInactiveDays
}

func (s *Service) calculateEnergy(
	currentEnergy int,
	inactiveDays int,
) int {

	energy := currentEnergy

	switch {

	case inactiveDays >= 7:
		energy -= 20

	case inactiveDays >= 3:
		energy -= 10

	case inactiveDays >= 1:
		energy -= 5
	}

	if energy < 10 {
		energy = 10
	}

	if energy > 100 {
		energy = 100
	}

	return energy
}

func (s *Service) determineKaiState(
	globalStreak int,
	inactiveDays int,
) string {

	switch {

	case globalStreak >= 7:
		return KaiStateProud

	case inactiveDays >= 7:
		return KaiStateSad

	case inactiveDays >= 1:
		return KaiStateSleeping

	default:
		return KaiStateCurious
	}
}

func (s *Service) determineKaiMode(
	inactiveDays int,
) string {

	switch {

	case inactiveDays >= 7:
		return KaiModeRecovery

	case inactiveDays >= 1:
		return KaiModeCompanion

	default:
		return KaiModeSoft
	}
}
