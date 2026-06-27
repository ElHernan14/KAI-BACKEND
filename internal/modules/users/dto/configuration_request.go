package usersdto

import (
	"strings"
	"time"
)

type UpdateUserConfigurationRequest struct {
	NotificacionesActivas       *bool   `json:"notificaciones_activas"`
	SonidosActivos              *bool   `json:"sonidos_activos"`
	MostrarRachas               *bool   `json:"mostrar_rachas"`
	ModoDiscreto                *bool   `json:"modo_discreto"`
	IntensidadKai               *string `json:"intensidad_kai" validate:"omitempty,oneof=baja normal alta"`
	HorarioRecordatorio         *string `json:"horario_recordatorio"`
	BloquearConPIN              *bool   `json:"bloquear_con_pin"`
	PermitirMensajesEmocionales *bool   `json:"permitir_mensajes_emocionales"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}
