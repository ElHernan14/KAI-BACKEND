package habitsmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserHabit struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"column:usuario_id;type:uuid;not null;index"`
	HabitCatalogID uuid.UUID `gorm:"column:habito_catalogo_id;type:uuid;not null"`
	Custom         bool      `gorm:"column:personalizado;default:false"`
	Active         bool      `gorm:"column:activo;default:true"`
	StartDate      time.Time `gorm:"column:fecha_inicio;type:date"`
	Personalized   bool      `gorm:"column:personalizado;default:false"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relaciones
	HabitCatalog HabitCatalog `gorm:"foreignKey:HabitCatalogID"`

	HabitRecords []HabitRecord `gorm:"foreignKey:UserHabitID"`
}

func (UserHabit) TableName() string {
	return "habitos_usuario"
}
