package messagesmodel

import (
	"time"

	"github.com/google/uuid"
)

type KaiMessage struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Type       string    `gorm:"column:tipo;type:varchar(50);not null"`
	Context    *string   `gorm:"column:contexto;type:varchar(50)"`
	Tone       *string   `gorm:"column:tono;type:varchar(30)"`
	Message    string    `gorm:"column:mensaje;type:text;not null"`
	Rarity     string    `gorm:"column:rareza;type:varchar(20);default:comun"`
	UnlockedBy *string   `gorm:"column:desbloqueado_por;type:varchar(100)"`
	Active     bool      `gorm:"column:activo;default:true"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (KaiMessage) TableName() string {
	return "mensajes_kai"
}
