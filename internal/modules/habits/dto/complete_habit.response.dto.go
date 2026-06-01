package habitsdto

import "github.com/google/uuid"

type CompleteHabitResponse struct {
	HabitoUsuarioID uuid.UUID `json:"habito_usuario_id"`

	XPGanada int `json:"xp_ganada"`

	RachaActual int `json:"racha_actual"`

	EnergiaActual int `json:"energia_actual"`

	NivelVinculo int `json:"nivel_vinculo"`

	AtributoDominanteID *uuid.UUID `json:"atributo_dominante_id"`
}
