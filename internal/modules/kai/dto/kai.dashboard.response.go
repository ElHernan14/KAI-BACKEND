package kaidto

import (
	"time"

	xpSummary "kai-back/internal/modules/xp/dto"
)

type KaiDashboardResponse struct {
	EstadoKai          KaiStateResponse                   `json:"estado_kai"`
	MensajeEmocional   string                             `json:"mensaje_emocional"`
	EventoEvolucion    EvolutionEventResponse             `json:"evento_evolucion"`
	Atributos          []KaiAttributeResponse             `json:"atributos"`
	CategoriaDominante *xpSummary.CategorySummaryResponse `json:"categoria_dominante"`
	CategoriaMenor     *xpSummary.CategorySummaryResponse `json:"categoria_menos_dominante"`
	ProgresoDiario     KaiProgressResponse                `json:"progreso_diario"`
}

type KaiStateResponse struct {
	EstadoActual      string     `json:"estado_actual"`
	EtapaActual       string     `json:"etapa_actual"`
	Energia           int        `json:"energia"`
	ImagenKai         *string    `json:"imagen_kai"`
	NivelVinculo      int        `json:"nivel_vinculo"`
	ModoRecuperacion  bool       `json:"modo_recuperacion"`
	ModoActual        *string    `json:"modo_actual"`
	UltimaEvolucion   *time.Time `json:"ultima_evolucion"`
	AtributoDominante string     `json:"atributo_dominante"`
}

type EvolutionEventResponse struct {
	Activo     bool       `json:"activo"`
	Etapa      string     `json:"etapa,omitempty"`
	IniciadoEn *time.Time `json:"iniciado_en,omitempty"`
	ExpiraEn   *time.Time `json:"expira_en,omitempty"`
}

type KaiAttributeResponse struct {
	Atributo string `json:"atributo"`
	Valor    int    `json:"valor"`
}

type KaiProgressResponse struct {
	Total       int `json:"total"`
	Completados int `json:"completados"`
	Pendientes  int `json:"pendientes"`
}
