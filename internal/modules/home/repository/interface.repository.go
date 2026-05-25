package home

import (
	"context"

	"github.com/google/uuid"
)

type HomeRepository interface {
	FindKaiSummary(ctx context.Context, userID uuid.UUID) (*KaiSummaryRow, error)
	FindTotalXP(ctx context.Context, userID uuid.UUID) (int, error)
	FindCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error)
	FindDailyHabits(ctx context.Context, userID uuid.UUID) ([]DailyHabitRow, error)
	FindFallbackMessage(ctx context.Context) (*string, error)
}
