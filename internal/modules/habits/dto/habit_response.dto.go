package habitsdto

import "time"

type HabitsViewResponse struct {
	DailyProgress DailyProgress       `json:"progreso_diario"`
	Habits        []UserHabitResponse `json:"habitos_usuario"`
}

type DailyProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completados"`
	Pending   int `json:"pendientes"`
}

type HabitRecordResponse struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"fecha"`
	Completed       bool      `json:"completado"`
	RegisteredValue *string   `json:"valor_registrado"`
	XPEarned        int       `json:"xp_ganada"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserHabitResponse struct {
	ID             string    `json:"habito_usuario_id"`
	HabitCatalogID string    `json:"habito_catalogo_id"`
	Name           string    `json:"nombre"`
	Description    *string   `json:"descripcion"`
	Category       *string   `json:"categoria"`
	CareType       *string   `json:"tipo_cuidado"`
	Difficulty     *string   `json:"dificultad"`
	BaseXP         int       `json:"xp_base"`
	Custom         bool      `json:"personalizado"`
	Active         bool      `json:"activo"`
	StartDate      time.Time `json:"fecha_inicio"`
	HabitImage     *string   `json:"imagen_habito"`
	CompletedToday bool      `json:"completado_hoy"`

	// Records []HabitRecordResponse `json:"registros"`
}
