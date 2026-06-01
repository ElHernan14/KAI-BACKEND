package xpmodel

import (
	kaimodels "kai-back/internal/modules/kai/models"

	"github.com/google/uuid"
)

type XPAttribute struct {
	ID uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`

	CategoryID uuid.UUID `gorm:"column:categoria_xp_id"`

	AttributeTypeID uuid.UUID `gorm:"column:atributo_kai_id"`

	Multiplier float64 `gorm:"column:multiplicador"`

	AttributeType *kaimodels.KaiAttributeType `gorm:"foreignKey:AttributeTypeID"`
}

func (XPAttribute) TableName() string {
	return "xp_atributos"
}
