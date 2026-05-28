package habitsdto

import "github.com/google/uuid"

type SelectHabitRequest struct {
	HabitCatalogID uuid.UUID `json:"habito_catalogo_id" validate:"required"`
}
