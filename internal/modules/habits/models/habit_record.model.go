package habitsmodel

import (
	"time"

	"github.com/google/uuid"
)

type HabitRecord struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UsuarioID       uuid.UUID `gorm:"column:usuario_id;type:uuid;not null;index"`
	UserHabitID     uuid.UUID `gorm:"column:habito_usuario_id;type:uuid;not null;index"`
	Fecha           time.Time `gorm:"column:fecha;type:date;not null"`
	Completado      bool      `gorm:"column:completado;default:false"`
	ValorRegistrado *string   `gorm:"column:valor_registrado;type:text"`
	XPGanada        int       `gorm:"column:xp_ganada;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relaciones
	UserHabit *UserHabit `gorm:"foreignKey:UserHabitID"`
}

func (HabitRecord) TableName() string {
	return "registros_habito"
}
