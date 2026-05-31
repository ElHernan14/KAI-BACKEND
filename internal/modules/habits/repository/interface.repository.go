package habitsRepository

import (
	"context"
	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"

	"github.com/google/uuid"
)

type HabitsRepository interface {
	FindUserHabits(ctx context.Context, userID uuid.UUID) ([]habitsmodel.UserHabit, error)
	CountDailyCompleted(ctx context.Context, userID uuid.UUID) (int, error)
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

	DeactivateHabit(
		ctx context.Context,
		userID uuid.UUID,
		habitUserID uuid.UUID,
	) error
}
