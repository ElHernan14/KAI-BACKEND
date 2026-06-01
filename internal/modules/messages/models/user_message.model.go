package messagesmodel

import (
	"time"

	"github.com/google/uuid"
)

type UserMessage struct {
	ID uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"column:usuario_id;type:uuid;not null"`

	MessageID uuid.UUID `gorm:"column:mensaje_kai_id;type:uuid;not null"`

	Read bool `gorm:"column:leido;default:false"`

	ShownAt time.Time `gorm:"column:mostrado_en"`

	CreatedAt time.Time `gorm:"column:created_at"`

	Message *KaiMessage `gorm:"foreignKey:MessageID"`
}

func (UserMessage) TableName() string {
	return "mensajes_usuario"
}
