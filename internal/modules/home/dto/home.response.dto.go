package homedto

import "time"

type HomeResponse struct {
	Kai            KaiHomeSummary         `json:"estado_kai"`
	Message        string                 `json:"mensaje_motivacional"`
	EvolutionEvent EvolutionEventResponse `json:"evento_evolucion"`
	TotalXP        int                    `json:"xp_total"`
	CurrentStreak  int                    `json:"racha_actual"`
	HabitsToday    []DailyHabitSummary    `json:"habitos_diarios"`
	DailyProgress  DailyProgress          `json:"progreso_diario"`
}

type KaiHomeSummary struct {
	CurrentState string  `json:"estado_actual"`
	CurrentStage string  `json:"etapa_actual"`
	Energy       int     `json:"energia"`
	BondLevel    int     `json:"nivel_vinculo"`
	RecoveryMode bool    `json:"modo_recuperacion"`
	KaiImage     *string `json:"imagen_kai"`
}

type EvolutionEventResponse struct {
	Active    bool       `json:"activo"`
	Stage     string     `json:"etapa,omitempty"`
	StartedAt *time.Time `json:"iniciado_en,omitempty"`
	ExpiresAt *time.Time `json:"expira_en,omitempty"`
	Message   string     `json:"mensaje,omitempty"`
}

type DailyHabitSummary struct {
	ID        string  `json:"habito_usuario_id"`
	Name      string  `json:"nombre"`
	Category  *string `json:"categoria"`
	Completed bool    `json:"completado"`
}

type DailyProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completados"`
	Pending   int `json:"pendientes"`
}
