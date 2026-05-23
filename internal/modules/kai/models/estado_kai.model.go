package kaimodel

import (
	"time"

	"github.com/google/uuid"
)

type KaiState struct {
	ID                  uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID              uuid.UUID  `gorm:"column:usuario_id;type:uuid;not null;unique" json:"usuario_id"`
	CurrentState        string     `gorm:"column:estado_actual;type:varchar(30);not null" json:"estado_actual"`
	CurrentStage        string     `gorm:"column:etapa_actual;type:varchar(30);not null" json:"etapa_actual"`
	Energy              int        `gorm:"column:energia;default:100" json:"energia"`
	KaiImage            *string    `gorm:"column:imagen_kai;type:varchar(100)" json:"imagen_kai,omitempty"`
	LastMessage         *string    `gorm:"column:ultimo_mensaje;type:text" json:"ultimo_mensaje,omitempty"`
	LastInteraction     *time.Time `gorm:"column:ultima_interaccion" json:"ultima_interaccion,omitempty"`
	DaysWithoutActivity int        `gorm:"column:dias_sin_actividad;default:0" json:"dias_sin_actividad"`
	RecoveryMode        bool       `gorm:"column:modo_recuperacion;default:false" json:"modo_recuperacion"`
	CreatedAt           time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CurrentMode         *string    `gorm:"column:modo_actual;type:varchar(30)" json:"modo_actual,omitempty"`
	DominantAttributeID *uuid.UUID `gorm:"column:atributo_dominante_id;type:uuid" json:"atributo_dominante_id,omitempty"`
	BondLevel           int        `gorm:"column:nivel_vinculo;default:1" json:"nivel_vinculo"`
	LastEvolution       *time.Time `gorm:"column:ultima_evolucion" json:"ultima_evolucion,omitempty"`

	DominantAttribute *KaiAttributeType `gorm:"foreignKey:DominantAttributeID"`
}

func (KaiState) TableName() string {
	return "estado_kai"
}
