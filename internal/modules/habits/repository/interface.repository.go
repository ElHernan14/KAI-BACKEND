package habitsRepository

import (
	"context"
	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HabitsRepository interface {
	FindUserHabits(ctx context.Context, userID uuid.UUID) ([]habitsmodel.UserHabit, error)
	CountDailyCompleted(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (int, error)
	FindCategories(
		ctx context.Context,
	) ([]habitsdto.HabitCategoryResponse, error)
	FindCatalogByCategoryExcludingUser(
		ctx context.Context,
		userID uuid.UUID,
		categoryID uuid.UUID,
	) ([]habitsmodel.HabitCatalog, error)

	FindCatalogByID(
		ctx context.Context,
		catalogID uuid.UUID,
	) (*habitsmodel.HabitCatalog, error)

	UserHasActiveHabit(
		ctx context.Context,
		userID uuid.UUID,
		catalogID uuid.UUID,
	) (bool, error)

	FindUserHabitByCatalogID(
		ctx context.Context,
		userID uuid.UUID,
		habitCatalogID uuid.UUID,
	) (*habitsmodel.UserHabit, error)

	ReactivateHabit(
		ctx context.Context,
		habitID uuid.UUID,
	) error

	ReactivateHabitWithTodayRecord(
		ctx context.Context,
		userHabit *habitsmodel.UserHabit,
	) error

	SelectHabit(
		ctx context.Context,
		userHabit *habitsmodel.UserHabit,
		initialRecord *habitsmodel.HabitRecord,
	) error

	FindHabitDetailByID(
		ctx context.Context,
		userID uuid.UUID,
		habitUserID uuid.UUID,
	) (*habitsmodel.UserHabit, error)

	FindUserHabitByID(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		habitID uuid.UUID,
	) (*habitsmodel.UserHabit, error)

	DeactivateHabit(
		ctx context.Context,
		habitID uuid.UUID,
	) error

	FindTodayRecord(
		ctx context.Context,
		tx *gorm.DB,
		userHabitID uuid.UUID,
		date time.Time,
	) (*habitsmodel.HabitRecord, error)

	UpdateHabitRecord(
		ctx context.Context,
		tx *gorm.DB,
		record *habitsmodel.HabitRecord,
	) error

	// FindStreakByUserHabitID(
	// 	ctx context.Context,
	// 	userHabitID uuid.UUID,
	// ) (*habitsmodel.Streak, error)

	CreateStreak(
		ctx context.Context,
		tx *gorm.DB,
		streak *habitsmodel.Streak,
	) error

	UpdateStreak(
		ctx context.Context,
		tx *gorm.DB,
		streak *habitsmodel.Streak,
	) error

	FindActiveUserHabits(
		ctx context.Context,
		userID uuid.UUID,
	) ([]habitsmodel.UserHabit, error)

	FindTodayRecords(
		ctx context.Context,
		userID uuid.UUID,
		date time.Time,
	) ([]habitsmodel.HabitRecord, error)

	CreateHabitRecords(
		ctx context.Context,
		tx *gorm.DB,
		records []habitsmodel.HabitRecord,
	) error

	FindStreakByHabit(
		ctx context.Context,
		tx *gorm.DB,
		userHabitID uuid.UUID,
	) (*habitsmodel.Streak, error)

	FindTodayHabitRecords(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		date time.Time,
	) ([]habitsmodel.HabitRecord, error)
}
