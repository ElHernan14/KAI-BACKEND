package kai

import (
	"context"
	"time"

	"github.com/google/uuid"

	habitsmodel "kai-back/internal/modules/habits/models"
	habitsRepository "kai-back/internal/modules/habits/repository"
	habitsServ "kai-back/internal/modules/habits/service"
	kaidto "kai-back/internal/modules/kai/dto"
	kaievolution "kai-back/internal/modules/kai/evolution"
	kaimodel "kai-back/internal/modules/kai/models"
	kaiRepository "kai-back/internal/modules/kai/repository"
	messagesmodel "kai-back/internal/modules/messages/models"
	messageRepository "kai-back/internal/modules/messages/repository"
	xpSummary "kai-back/internal/modules/xp/dto"
	xpmodel "kai-back/internal/modules/xp/models"
	xpRepository "kai-back/internal/modules/xp/repository"
	userActivitySynchronizationService "kai-back/internal/services/user_activity_synchronization"
)

type ServicePort interface {
	GetKaiDashboard(
		ctx context.Context,
		userID uuid.UUID,
	) (*kaidto.KaiDashboardResponse, error)
}

type KaiService struct {
	kaiRepository                      kaiRepository.KaiRepository
	messageRepository                  messageRepository.MessageRepositoryPort
	habitsRepository                   habitsRepository.HabitsRepository
	xpRepository                       xpRepository.XpRepository
	habitsDailyRecordsService          habitsServ.HabitsDailyRecordsServicePort
	UserActivitySynchronizationService userActivitySynchronizationService.ServicePort
}

func NewKaiService(
	kaiRepository kaiRepository.KaiRepository,
	messageRepository messageRepository.MessageRepositoryPort,
	habitsRepository habitsRepository.HabitsRepository,
	xpRepository xpRepository.XpRepository,
	habitsDailyRecordsService habitsServ.HabitsDailyRecordsServicePort,
	userActivitySyncService userActivitySynchronizationService.ServicePort,
) *KaiService {
	return &KaiService{
		kaiRepository:                      kaiRepository,
		messageRepository:                  messageRepository,
		habitsRepository:                   habitsRepository,
		xpRepository:                       xpRepository,
		habitsDailyRecordsService:          habitsDailyRecordsService,
		UserActivitySynchronizationService: userActivitySyncService,
	}
}

func (s *KaiService) GetKaiDashboard(
	ctx context.Context,
	userID uuid.UUID,
) (*kaidto.KaiDashboardResponse, error) {

	// Asegurar que los registros diarios de hábitos estén creados para hoy
	err := s.habitsDailyRecordsService.
		EnsureTodayHabitRecords(
			ctx,
			userID,
		)
	if err != nil {
		return nil, err
	}

	// Aseguramos Sincronizar toda la información temporal del usuario dependiente del paso del tiempo y de su actividad reciente.
	err = s.UserActivitySynchronizationService.SyncUserActivityState(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// Obtener datos necesarios para el dashboard
	state, err := s.kaiRepository.FindKaiStateByUserID(
		ctx,
		nil,
		userID,
	)
	if err != nil {
		return nil, err
	}

	attributes, err := s.kaiRepository.FindUserKaiAttributes(
		ctx,
		nil,
		userID,
	)
	if err != nil {
		return nil, err
	}

	records, err := s.habitsRepository.FindTodayHabitRecords(
		ctx,
		nil,
		userID,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	message, err := s.messageRepository.FindLastUserMessage(
		ctx,
		nil,
		userID,
	)
	if err != nil {
		return nil, err
	}

	userXP, err := s.xpRepository.FindUserXPByUserID(
		ctx,
		nil,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// construir datos orientados a vista
	motivationalMessage := s.buildMessage(
		message,
	)

	dailyProgress := s.buildDailyProgress(
		records,
	)

	dominantCategory := s.buildDominantCategory(
		userXP,
	)

	weakestCategory := s.buildWeakestCategory(
		userXP,
	)

	evolutionActive, evolutionExpiresAt := kaievolution.EventWindow(state.LastEvolution, time.Now())

	// - mapear DTO
	return &kaidto.KaiDashboardResponse{
		EstadoKai: kaidto.KaiStateResponse{
			EstadoActual:      state.CurrentState,
			EtapaActual:       state.CurrentStage,
			ModoActual:        state.CurrentMode,
			Energia:           state.Energy,
			NivelVinculo:      state.BondLevel,
			ModoRecuperacion:  state.RecoveryMode,
			ImagenKai:         state.KaiImage,
			UltimaEvolucion:   state.LastEvolution,
			AtributoDominante: state.DominantAttribute.Name,
		},

		Atributos: s.mapKaiAttributes(attributes),

		MensajeEmocional: motivationalMessage,
		EventoEvolucion: kaidto.EvolutionEventResponse{
			Activo:     evolutionActive,
			Etapa:      state.CurrentStage,
			IniciadoEn: state.LastEvolution,
			ExpiraEn:   evolutionExpiresAt,
		},

		ProgresoDiario: dailyProgress,

		CategoriaDominante: dominantCategory,

		CategoriaMenor: weakestCategory,
	}, nil
}

func (s *KaiService) mapKaiAttributes(
	attrs []kaimodel.KaiAttribute,
) []kaidto.KaiAttributeResponse {

	result := make(
		[]kaidto.KaiAttributeResponse,
		0,
		len(attrs),
	)

	for _, attr := range attrs {

		var name string

		if attr.AttributeType != nil {
			name = attr.AttributeType.Name
		}

		result = append(
			result,
			kaidto.KaiAttributeResponse{
				Atributo: name,
				Valor:    attr.Value,
			},
		)
	}

	return result
}

func (s *KaiService) buildDailyProgress(
	records []habitsmodel.HabitRecord,
) kaidto.KaiProgressResponse {

	total := len(records)
	completed := 0

	for _, record := range records {
		if record.Completado {
			completed++
		}
	}

	return kaidto.KaiProgressResponse{
		Total:       total,
		Completados: completed,
		Pendientes:  total - completed,
	}
}

func (s *KaiService) buildDominantCategory(
	userXP []xpmodel.UserXP,
) *xpSummary.CategorySummaryResponse {

	if len(userXP) == 0 {
		return nil
	}

	dominant := userXP[0]

	for _, xp := range userXP {
		if xp.Value > dominant.Value {
			dominant = xp
		}
	}

	if dominant.Category == nil {
		return nil
	}

	return &xpSummary.CategorySummaryResponse{
		Nombre: dominant.Category.Name,
		Valor:  dominant.Value,
	}
}

func (s *KaiService) buildWeakestCategory(
	userXP []xpmodel.UserXP,
) *xpSummary.CategorySummaryResponse {

	if len(userXP) == 0 {
		return nil
	}

	lowest := userXP[0]

	for _, xp := range userXP {
		if xp.Value < lowest.Value {
			lowest = xp
		}
	}

	if lowest.Category == nil {
		return nil
	}

	return &xpSummary.CategorySummaryResponse{
		Nombre: lowest.Category.Name,
		Valor:  lowest.Value,
	}
}

func (s *KaiService) buildMessage(
	message *messagesmodel.UserMessage,
) string {

	if message == nil || message.KaiMessage == nil {
		return ""
	}

	return message.KaiMessage.Message
}
