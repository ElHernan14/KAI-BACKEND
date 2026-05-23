package kaimodel

import (
	"time"

	"github.com/google/uuid"
)

type KaiAttributeType struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"column:nombre"`
	Description *string   `gorm:"column:descripcion"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (KaiAttributeType) TableName() string {
	return "tipos_atributo_kai"
}
