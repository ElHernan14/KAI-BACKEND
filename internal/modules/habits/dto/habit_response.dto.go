package habitsdto

import "time"

type HabitRecordResponse struct {
	Fecha           time.Time `json:"fecha"`
	Completado      bool      `json:"completado"`
	ValorRegistrado *string   `json:"valor_registrado"`
	XPGanada        int       `json:"xp_ganada"`
}

type UserHabitResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"nombre"`
	Description *string `json:"descripcion"`

	Category   *string `json:"categoria"`
	CareType   *string `json:"tipo_cuidado"`
	Difficulty *string `json:"dificultad"`

	BaseXP int `json:"xp_base"`

	Active bool `json:"activo"`

	HabitImage *string `json:"imagen_habito"`

	Records []HabitRecordResponse `json:"registros"`
}
