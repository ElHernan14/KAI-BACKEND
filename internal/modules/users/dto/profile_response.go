package usersdto

import "time"

type UserProfileResponse struct {
	Usuario              UserProfileDataResponse   `json:"usuario"`
	ConfiguracionUsuario UserConfigurationResponse `json:"configuracion_usuario"`
	Resumen              UserProfileSummary        `json:"resumen"`
}

type UserProfileDataResponse struct {
	ID            string    `json:"id"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email"`
	Username      *string   `json:"username"`
	FotoPerfil    *string   `json:"foto_perfil"`
	PerfilBase    *string   `json:"perfil_base"`
	EtapaKai      string    `json:"etapa_kai"`
	FechaRegistro time.Time `json:"fecha_registro"`
}

type UserConfigurationResponse struct {
	NotificacionesActivas       bool    `json:"notificaciones_activas"`
	SonidosActivos              bool    `json:"sonidos_activos"`
	MostrarRachas               bool    `json:"mostrar_rachas"`
	ModoDiscreto                bool    `json:"modo_discreto"`
	IntensidadKai               string  `json:"intensidad_kai"`
	HorarioRecordatorio         *string `json:"horario_recordatorio"`
	BloquearConPIN              bool    `json:"bloquear_con_pin"`
	PermitirMensajesEmocionales bool    `json:"permitir_mensajes_emocionales"`
}

type UserProfileSummary struct {
	XPTotal      int `json:"xp_total"`
	RachaGlobal  int `json:"racha_global"`
	DiasInactivo int `json:"dias_inactivo"`
}
