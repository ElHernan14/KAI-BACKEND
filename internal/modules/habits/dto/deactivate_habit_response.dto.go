package habitsdto

import "github.com/google/uuid"

type DeactivateHabitResponse struct {
	HabitID uuid.UUID `json:"habit_id"`
	Active  bool      `json:"active"`
	Message string    `json:"message"`
}
