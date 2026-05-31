package habitsdto

import (
	"time"

	"github.com/google/uuid"
)

type HabitDetailResponse struct {
	Habito       UserHabitResponse     `json:"habito"`
	Estadisticas HabitStatsResponse    `json:"estadisticas"`
	Registros    []HabitRecordResponse `json:"registros"`
}

type HabitStatsResponse struct {
	TotalCompletados int `json:"total_completados"`
	XPTotal          int `json:"xp_total"`
}

type HabitRecordResponse struct {
	ID              string    `json:"id"`
	UserHabitID     uuid.UUID `json:"habito_usuario_id"`
	Date            time.Time `json:"fecha"`
	Completed       bool      `json:"completado"`
	RegisteredValue *string   `json:"valor_registrado"`
	XPEarned        int       `json:"xp_ganada"`
	CreatedAt       time.Time `json:"created_at"`
}
