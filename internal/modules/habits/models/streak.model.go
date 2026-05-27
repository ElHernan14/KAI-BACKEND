package habitsmodel

import (
	"time"

	"github.com/google/uuid"
)

type Streak struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID  `gorm:"column:usuario_id;type:uuid;not null;index"`
	UserHabitID    uuid.UUID  `gorm:"column:habito_usuario_id;type:uuid;not null;index"`
	CurrentDays    int        `gorm:"column:dias_actuales;default:0"`
	HistoricalBest int        `gorm:"column:record_historico;default:0"`
	StartDate      *time.Time `gorm:"column:inicio_racha;type:date"`
	LastActivity   *time.Time `gorm:"column:ultima_actividad;type:date"`
	Active         bool       `gorm:"column:activa;default:true"`
	Protected      bool       `gorm:"column:protegida;default:false"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Streak) TableName() string {
	return "rachas"
}
