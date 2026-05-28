package habitsdto

import "github.com/google/uuid"

type HabitCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"nombre"`
	Description *string   `json:"descripcion"`
}

type HabitCatalogResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"nombre"`
	Description  *string   `json:"descripcion"`
	Category     *string   `json:"categoria"`
	CareType     *string   `json:"tipo_cuidado"`
	Difficulty   *string   `json:"dificultad"`
	BaseXP       int       `json:"xp_base"`
	Premium      bool      `json:"es_premium"`
	HabitImage   *string   `json:"imagen_habito"`
	CategoryXPID uuid.UUID `json:"categoria_xp_id"`
}

type SelectHabitResponse struct {
	ID             uuid.UUID `json:"habito_usuario_id"`
	HabitCatalogID uuid.UUID `json:"habito_catalogo_id"`
	Active         bool      `json:"activo"`
}
