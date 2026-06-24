package habitcompletion

import (
	"context"

	habitsdto "kai-back/internal/modules/habits/dto"

	"github.com/google/uuid"
)

type ServicePort interface {
	CompleteHabit(
		ctx context.Context,
		userID uuid.UUID,
		userHabitID uuid.UUID,
		value *string,
	) (*habitsdto.CompleteHabitResponse, error)
}
