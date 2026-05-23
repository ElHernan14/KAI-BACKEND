package kaidto

import (
	"time"

	"github.com/google/uuid"
)

type KaiStateSummary struct {
	CurrentState      string                     `json:"estado_actual"`
	CurrentStage      string                     `json:"etapa_actual"`
	CurrentMode       string                     `json:"modo_actual"`
	Energy            int                        `json:"energia"`
	BondLevel         int                        `json:"nivel_vinculo"`
	RecoveryMode      bool                       `json:"modo_recuperacion"`
	DominantAttribute *DominantAttributeResponse `json:"atributoDominante"`

	LastEvolution *time.Time `json:"ultima_evolucion,omitempty"`
}

type KaiAttributeSummary struct {
	AttributeID   string `json:"atributo_kai_id"`
	AttributeName string `json:"nombre"`
	Value         int    `json:"valor"`
}

type DominantAttributeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"nombre"`
}
