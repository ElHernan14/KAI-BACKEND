package kaidto

type KaiStateSummary struct {
	CurrentState string `json:"estado_actual"`
	CurrentStage string `json:"etapa_actual"`
	Energy       int    `json:"energia"`
	BondLevel    int    `json:"nivel_vinculo"`
	RecoveryMode bool   `json:"modo_recuperacion"`
}
