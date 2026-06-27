package usersmodel

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserConfiguration struct {
	ID                     uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID                 uuid.UUID      `gorm:"column:usuario_id;type:uuid;not null;unique" json:"usuario_id"`
	NotificationsEnabled   bool           `gorm:"column:notificaciones_activas;default:true" json:"notificaciones_activas"`
	SoundsEnabled          bool           `gorm:"column:sonidos_activos;default:true" json:"sonidos_activos"`
	ShowStreaks            bool           `gorm:"column:mostrar_rachas;default:true" json:"mostrar_rachas"`
	DiscreteMode           bool           `gorm:"column:modo_discreto;default:false" json:"modo_discreto"`
	KaiIntensity           string         `gorm:"column:intensidad_kai;type:varchar(20);default:normal" json:"intensidad_kai"`
	ReminderTime           sql.NullString `gorm:"column:horario_recordatorio;type:time" json:"horario_recordatorio,omitempty"`
	LockWithPIN            bool           `gorm:"column:bloquear_con_pin;default:false" json:"bloquear_con_pin"`
	AllowEmotionalMessages bool           `gorm:"column:permitir_mensajes_emocionales;default:true" json:"permitir_mensajes_emocionales"`
	CreatedAt              time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (UserConfiguration) TableName() string {
	return "configuracion_usuario"
}
