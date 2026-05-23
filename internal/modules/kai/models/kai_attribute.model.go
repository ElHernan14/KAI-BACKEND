package kaimodel

import (
	"time"

	"github.com/google/uuid"
)

type KaiAttribute struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID          uuid.UUID `gorm:"column:usuario_id"`
	AttributeTypeID uuid.UUID `gorm:"column:atributo_kai_id"`

	Value int `gorm:"column:valor"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	AttributeType *KaiAttributeType `gorm:"foreignKey:AttributeTypeID"`
}

func (KaiAttribute) TableName() string {
	return "atributos_kai"
}
