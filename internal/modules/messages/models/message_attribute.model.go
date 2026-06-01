package messagesmodel

import (
	kaimodel "kai-back/internal/modules/kai/models"

	"github.com/google/uuid"
)

type MessageAttribute struct {
	ID uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`

	MessageID uuid.UUID `gorm:"column:mensaje_kai_id;type:uuid;not null"`

	AttributeTypeID uuid.UUID `gorm:"column:atributo_kai_id;type:uuid;not null"`

	MinimumLevel int `gorm:"column:nivel_minimo;default:1"`

	Message *KaiMessage `gorm:"foreignKey:MessageID"`

	AttributeType *kaimodel.KaiAttributeType `gorm:"foreignKey:AttributeTypeID"`
}

func (MessageAttribute) TableName() string {
	return "mensaje_atributos"
}
