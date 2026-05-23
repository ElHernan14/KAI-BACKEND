package xpmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserXP struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"column:usuario_id"`

	CategoryID uuid.UUID `gorm:"column:categoria_xp_id"`

	Value int `gorm:"column:valor"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	Category *XPCategory `gorm:"foreignKey:CategoryID"`
}

func (UserXP) TableName() string {
	return "xp_usuario"
}
