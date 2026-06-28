package habitcompletion

import (
	"context"
	"errors"
	habitsmodel "kai-back/internal/modules/habits/models"
	kaievolution "kai-back/internal/modules/kai/evolution"
	kaimodel "kai-back/internal/modules/kai/models"
	messagesmodel "kai-back/internal/modules/messages/models"
	errorHandler "kai-back/internal/shared/errors"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) validateHabit(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	userHabitID uuid.UUID,
) (*habitsmodel.UserHabit,
	*habitsmodel.HabitRecord,
	error) {

	habit, err := s.habitRepo.FindUserHabitByID(
		ctx,
		tx,
		userID,
		userHabitID,
	)
	if err != nil {
		return nil, nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"habito no encontrado",
		)
	}

	if !habit.Active {
		return nil, nil, errorHandler.NewAppError(
			http.StatusConflict,
			"el habito se encuentra inactivo",
		)
	}

	record, err := s.habitRepo.FindTodayRecord(
		ctx,
		tx,
		userHabitID,
		time.Now(),
	)
	if err != nil {
		return nil, nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"registro del dia no encontrado",
		)
	}

	if record.Completado {
		return nil, nil, errorHandler.NewAppError(
			http.StatusConflict,
			"el habito ya fue completado hoy",
		)
	}

	return habit, record, nil
}

func (s *Service) completeRecord(
	ctx context.Context,
	tx *gorm.DB,
	habit *habitsmodel.UserHabit,
	value *string,
	record *habitsmodel.HabitRecord,
) (*habitsmodel.HabitRecord, error) {

	if record.Completado {
		return nil, errorHandler.NewAppError(
			http.StatusConflict,
			"el habito ya fue completado hoy",
		)
	}

	record.Completado = true
	record.ValorRegistrado = value
	record.XPGanada = habit.HabitCatalog.BaseXP

	err := s.habitRepo.UpdateHabitRecord(
		ctx,
		tx,
		record,
	)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al actualizar el registro del habito",
		)
	}

	return record, nil
}

func (s *Service) updateStreak(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	habit *habitsmodel.UserHabit,
) (*habitsmodel.Streak, error) {

	streak, err := s.habitRepo.FindStreakByHabit(
		ctx,
		tx,
		habit.ID,
	)

	// =========================
	// Primera vez
	// =========================

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// crear racha
		today := time.Now()

		newStreak := &habitsmodel.Streak{
			UserID:         userID,
			UserHabitID:    habit.ID,
			CurrentDays:    1,
			HistoricalBest: 1,
			StartDate:      &today,
			LastActivity:   &today,
			Active:         true,
			Protected:      false,
		}

		err = s.habitRepo.CreateStreak(
			ctx,
			tx,
			newStreak,
		)
		if err != nil {
			return nil, errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al crear la racha del habito",
			)
		}

		return newStreak, nil
	} else if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar la racha del habito",
		)
	}

	today := time.Now()

	// =========================
	// Ya actualizada hoy
	// =========================

	if streak.LastActivity != nil &&
		sameDay(
			*streak.LastActivity,
			today,
		) {

		return streak, nil
	}

	// =========================
	// Continúa racha
	// =========================

	if streak.LastActivity != nil &&
		isYesterday(
			*streak.LastActivity,
			today,
		) {

		streak.CurrentDays++

		if streak.CurrentDays >
			streak.HistoricalBest {

			streak.HistoricalBest =
				streak.CurrentDays
		}

		streak.LastActivity = &today

		err = s.habitRepo.UpdateStreak(
			ctx,
			tx,
			streak,
		)
		if err != nil {
			return nil, errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al actualizar la racha del habito",
			)
		}

		return streak, nil
	}

	// =========================
	// Racha rota
	// =========================

	streak.CurrentDays = 1
	streak.StartDate = &today
	streak.LastActivity = &today
	streak.Active = true

	err = s.habitRepo.UpdateStreak(
		ctx,
		tx,
		streak,
	)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al actualizar la racha del habito",
		)
	}

	return streak, nil
}

func (s *Service) grantXP(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	habit *habitsmodel.UserHabit,
) (int, error) {

	xpGranted := habit.
		HabitCatalog.
		BaseXP

	userXP, err := s.xpRepo.
		FindUserXPByCategory(
			ctx,
			tx,
			userID,
			habit.HabitCatalog.XPCategoryID,
		)

	if err != nil {
		return 0, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar el XP del usuario",
		)
	}

	userXP.Value += xpGranted

	err = s.xpRepo.
		UpdateUserXP(
			ctx,
			tx,
			userXP,
		)
	if err != nil {
		return 0, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al actualizar el XP del usuario",
		)
	}

	return xpGranted, nil
}

func (s *Service) grantKaiAttributes(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	habit *habitsmodel.UserHabit,
	xpGranted int,
) error {

	mappings, err := s.xpRepo.
		FindXPAttributesByCategory(
			ctx,
			tx,
			habit.HabitCatalog.XPCategoryID,
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar los atributos asociados al XP",
		)
	}

	for _, mapping := range mappings {

		attr, err := s.kaiRepo.
			FindKaiAttribute(
				ctx,
				tx,
				userID,
				mapping.AttributeTypeID,
			)
		if err != nil {
			return errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al buscar el atributo de Kai",
			)
		}

		gain := int(
			math.Round(
				float64(xpGranted) *
					mapping.Multiplier,
			),
		)

		attr.Value += gain

		err = s.kaiRepo.
			UpdateKaiAttribute(
				ctx,
				tx,
				attr,
			)
		if err != nil {
			return errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al actualizar el atributo de Kai",
			)
		}
	}

	return nil
}

func (s *Service) recalculateDominantAttribute(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*kaimodel.KaiAttribute, error) {

	attributes, err := s.kaiRepo.
		FindUserKaiAttributes(
			ctx,
			tx,
			userID,
		)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar los atributos de Kai del usuario",
		)
	}

	if len(attributes) == 0 {
		return nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"no se encontraron atributos de Kai para el usuario",
		)
	}

	dominant := attributes[0]

	for _, attr := range attributes {

		if attr.Value > dominant.Value {
			dominant = attr
		}
	}

	kaiState, err := s.kaiRepo.
		FindKaiStateByUserID(
			ctx,
			tx,
			userID,
		)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar el estado de Kai del usuario",
		)
	}

	kaiState.DominantAttributeID =
		&dominant.AttributeTypeID

	err = s.kaiRepo.
		UpdateKaiState(
			ctx,
			tx,
			kaiState,
		)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al actualizar el estado de Kai del usuario",
		)
	}

	return &dominant, nil
}

type evolutionResult struct {
	Evolved        bool
	PreviousStage  string
	State          *kaimodel.KaiState
	EventActive    bool
	EventExpiresAt *time.Time
	Message        *string
}

func (s *Service) generateEvolutionMessage(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	evolution *evolutionResult,
) error {
	messageType := messagesmodel.MessageTypeEvolutionYoung
	if kaievolution.NormalizeStage(evolution.State.CurrentStage) == "adulto" {
		messageType = messagesmodel.MessageTypeEvolutionAdult
	}

	message, err := s.messageRepo.FindRandomMessageByType(ctx, tx, userID, messageType)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar un mensaje de evolucion",
		)
	}
	if message == nil {
		return nil
	}

	if err := s.messageRepo.CreateUserMessage(ctx, tx, &messagesmodel.UserMessage{
		UserID:    userID,
		MessageID: message.ID,
		Read:      false,
		ShownAt:   time.Now(),
	}); err != nil {
		return err
	}
	if err := s.kaiRepo.UpdateLastMessage(ctx, tx, userID, message.Message); err != nil {
		return err
	}

	evolution.Message = &message.Message
	return nil
}

func (s *Service) updateKaiState(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	habit *habitsmodel.UserHabit,
	streak *habitsmodel.Streak,
	dominantAttribute *kaimodel.KaiAttribute,
) (*evolutionResult, error) {
	kaiState, err := s.kaiRepo.FindKaiStateByUserID(ctx, tx, userID)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar el estado de Kai del usuario",
		)
	}

	now := time.Now()
	totalXP, err := s.xpRepo.FindTotalUserXP(ctx, tx, userID)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al calcular el XP total del usuario",
		)
	}

	completedToday, err := s.habitRepo.CountDailyCompleted(ctx, tx, userID)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al calcular los habitos completados del dia",
		)
	}

	currentStreak := 0
	if streak != nil {
		currentStreak = streak.CurrentDays
	}

	previousStage := kaievolution.NormalizeStage(kaiState.CurrentStage)
	nextStage := calculateStage(totalXP)
	stageChanged := previousStage != nextStage

	attributeName := ""
	if dominantAttribute != nil {
		kaiState.DominantAttributeID = &dominantAttribute.AttributeTypeID
		if dominantAttribute.AttributeType != nil {
			attributeName = dominantAttribute.AttributeType.Name
		}
	}

	kaiState.LastInteraction = &now
	kaiState.DaysWithoutActivity = 0
	kaiState.RecoveryMode = false
	kaiState.CurrentStage = nextStage
	kaiState.CurrentState = calculateState(completedToday, currentStreak)

	currentMode := calculateMode(attributeName, currentStreak, stageChanged)
	kaiState.CurrentMode = &currentMode
	kaiState.Energy = clampInt(
		kaiState.Energy+calculateEnergyGain(habit.HabitCatalog.Difficulty, currentStreak),
		kaiEnergyMin,
		kaiEnergyMax,
	)
	kaiState.BondLevel = clampInt(
		kaiState.BondLevel+calculateBondGain(currentStreak),
		kaiBondMin,
		kaiBondMax,
	)

	if stageChanged {
		kaiState.LastEvolution = &now
		imageKey := kaievolution.ImageKeyForStage(nextStage)
		kaiState.KaiImage = &imageKey
	}

	if err := s.kaiRepo.UpdateKaiState(ctx, tx, kaiState); err != nil {
		return nil, err
	}

	eventActive, eventExpiresAt := kaievolution.EventWindow(kaiState.LastEvolution, now)
	return &evolutionResult{
		Evolved:        stageChanged,
		PreviousStage:  previousStage,
		State:          kaiState,
		EventActive:    eventActive,
		EventExpiresAt: eventExpiresAt,
	}, nil
}

func (s *Service) generateMotivationalMessage(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	streak *habitsmodel.Streak,
	dominantAttribute *kaimodel.KaiAttribute,
) error {

	if dominantAttribute == nil {
		return nil
	}

	completedToday, err := s.habitRepo.
		CountDailyCompleted(
			ctx,
			tx,
			userID,
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al calcular contexto de mensajes de Kai",
		)
	}

	currentStreak := 0
	if streak != nil {
		currentStreak = streak.CurrentDays
	}

	attributeName := ""
	if dominantAttribute.AttributeType != nil {
		attributeName = dominantAttribute.AttributeType.Name
	}

	state := calculateState(
		completedToday,
		currentStreak,
	)
	mode := calculateMode(
		attributeName,
		currentStreak,
		false,
	)

	messageTypes, contexts := messageFiltersForRules(
		state,
		mode,
	)

	message, err := s.messageRepo.
		FindRandomMessageByRules(
			ctx,
			tx,
			dominantAttribute.AttributeTypeID,
			dominantAttribute.Value,
			messageTypes,
			contexts,
		)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al buscar un mensaje motivacional",
			)
		}

		message, err = s.messageRepo.
			FindRandomMessageByAttribute(
				ctx,
				tx,
				dominantAttribute.AttributeTypeID,
			)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			return errorHandler.NewAppError(
				http.StatusInternalServerError,
				"error al buscar un mensaje motivacional",
			)
		}
	}

	userMessage := &messagesmodel.UserMessage{
		UserID:    userID,
		MessageID: message.ID,
		Read:      false,
		ShownAt:   time.Now(),
	}

	err = s.messageRepo.
		CreateUserMessage(
			ctx,
			tx,
			userMessage,
		)
	if err != nil {
		return err
	}

	return s.kaiRepo.
		UpdateLastMessage(
			ctx,
			tx,
			userID,
			message.Message,
		)
}

func sameDay(
	a time.Time,
	b time.Time,
) bool {

	return a.Format("2006-01-02") ==
		b.Format("2006-01-02")
}

func isYesterday(
	last time.Time,
	today time.Time,
) bool {

	yesterday := today.AddDate(
		0,
		0,
		-1,
	)

	return sameDay(
		last,
		yesterday,
	)
}
