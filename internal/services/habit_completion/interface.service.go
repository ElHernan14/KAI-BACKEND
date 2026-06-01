package habitcompletion

import (
	"context"

	"github.com/google/uuid"
)

type ServicePort interface {
	CompleteHabit(
		ctx context.Context,
		userID uuid.UUID,
		userHabitID uuid.UUID,
		value *string,
	) error
}
