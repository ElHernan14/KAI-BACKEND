package habitsdto

import (
	"time"

	"github.com/google/uuid"
)

type CompleteHabitResponse struct {
	HabitoUsuarioID uuid.UUID `json:"habito_usuario_id"`
	XPGanada        int       `json:"xp_ganada"`
	RachaActual     int       `json:"racha_actual"`
	EnergiaActual   int       `json:"energia_actual"`
	NivelVinculo    int       `json:"nivel_vinculo"`

	AtributoDominanteID *uuid.UUID `json:"atributo_dominante_id"`

	Evoluciono      bool                   `json:"evoluciono"`
	EtapaAnterior   string                 `json:"etapa_anterior"`
	EtapaActual     string                 `json:"etapa_actual"`
	EventoEvolucion EvolutionEventResponse `json:"evento_evolucion"`
}

type EvolutionEventResponse struct {
	Activo     bool       `json:"activo"`
	Etapa      string     `json:"etapa,omitempty"`
	IniciadoEn *time.Time `json:"iniciado_en,omitempty"`
	ExpiraEn   *time.Time `json:"expira_en,omitempty"`
}
