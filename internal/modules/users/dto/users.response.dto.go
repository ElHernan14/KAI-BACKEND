package usersdto

import (
	kaiStateDto "kai-back/internal/modules/kai/dto"

	"github.com/google/uuid"
)

type MeResponse struct {
	ID            uuid.UUID                    `json:"id"`
	Name          string                       `json:"nombre"`
	Email         string                       `json:"email"`
	ProfileBase   *string                      `json:"perfil_base"`
	KaiStage      string                       `json:"etapa_kai"`
	GlobalStreak  int                          `json:"racha_global"`
	InactiveDays  int                          `json:"dias_inactivo"`
	KaiState      *kaiStateDto.KaiStateSummary `json:"estado_kai"`
	Configuration *UserConfigResponse          `json:"configuracion_usuario"`
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
