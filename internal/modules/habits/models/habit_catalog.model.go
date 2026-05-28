package habitsmodel

import (
	"time"

	xpmodel "kai-back/internal/modules/xp/models"

	"github.com/google/uuid"
)

type HabitCatalog struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"column:nombre;type:varchar(100);not null"`
	Description  *string   `gorm:"column:descripcion;type:text"`
	Category     *string   `gorm:"column:categoria;type:varchar(50)"`
	CareType     *string   `gorm:"column:tipo_cuidado;type:varchar(50)"`
	Difficulty   *string   `gorm:"column:dificultad;type:varchar(20)"`
	BaseXP       int       `gorm:"column:xp_base;default:10"`
	IsPremium    bool      `gorm:"column:es_premium;default:false"`
	Active       bool      `gorm:"column:activo;default:true"`
	HabitImage   *string   `gorm:"column:imagen_habito;type:varchar(100)"`
	XPCategoryID uuid.UUID `gorm:"column:categoria_xp_id;type:uuid"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relaciones
	XPCategory xpmodel.XPCategory `gorm:"foreignKey:XPCategoryID"`
}

func (HabitCatalog) TableName() string {
	return "habitos_catalogo"
}
