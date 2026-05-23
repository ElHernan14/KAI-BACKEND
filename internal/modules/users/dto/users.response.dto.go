package usersdto

import (
	habitsdto "kai-back/internal/modules/habits/dto"
	kaiStateDto "kai-back/internal/modules/kai/dto"
	kaidto "kai-back/internal/modules/kai/dto"
	xpdto "kai-back/internal/modules/xp/dto"
)

type MeResponse struct {
	User          *UserResponse                 `json:"usuario"`
	KaiState      *kaiStateDto.KaiStateSummary  `json:"estado_kai"`
	Configuration *UserConfigResponse           `json:"configuracion_usuario"`
	XP            []xpdto.UserXPSummary         `json:"xp_usuario"`
	Attributes    []kaidto.KaiAttributeSummary  `json:"atributos_kai"`
	Habits        []habitsdto.UserHabitResponse `json:"habitos"`
}

type UserResponse struct {
	Name         string  `json:"nombre"`
	Email        string  `json:"email"`
	ProfileBase  *string `json:"perfil_base"`
	KaiStage     string  `json:"etapa_kai"`
	GlobalStreak int     `json:"racha_global"`
	InactiveDays int     `json:"dias_inactivo"`
}

type UserConfigResponse struct {
	NotificationsEnabled   bool    `json:"notificaciones_activas"`
	SoundsEnabled          bool    `json:"sonidos_activos"`
	ShowStreaks            bool    `json:"mostrar_rachas"`
	DiscreteMode           bool    `json:"modo_discreto"`
	KaiIntensity           string  `json:"intensidad_kai"`
	ReminderTime           *string `json:"horario_recordatorio"`
	LockWithPIN            bool    `json:"bloquear_con_pin"`
	AllowEmotionalMessages bool    `json:"permitir_mensajes_emocionales"`
}
