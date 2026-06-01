package habitcompletion

import (
	"context"
	"errors"
	habitsmodel "kai-back/internal/modules/habits/models"
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
	userID uuid.UUID,
	userHabitID uuid.UUID,
) (*habitsmodel.UserHabit,
	*habitsmodel.HabitRecord,
	error) {

	habit, err := s.habitRepo.FindUserHabitByID(
		ctx,
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
) (*uuid.UUID, error) {

	attributes, err := s.kaiRepo.
		FindUserKaiAttributes(
			ctx,
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

	return kaiState.DominantAttributeID, nil
}

func (s *Service) updateKaiState(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	dominantAttributeID *uuid.UUID,
) error {

	kaiState, err := s.kaiRepo.
		FindKaiStateByUserID(
			ctx,
			userID,
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar el estado de Kai del usuario",
		)
	}

	now := time.Now()

	kaiState.LastInteraction = &now

	kaiState.DaysWithoutActivity = 0

	kaiState.RecoveryMode = false

	kaiState.DominantAttributeID =
		dominantAttributeID

	kaiState.Energy += 10

	if kaiState.Energy > 100 {
		kaiState.Energy = 100
	}

	kaiState.BondLevel += 1

	// TODO:
	// Recalcular estado_actual
	// Recalcular etapa_actual
	// Verificar evolución

	return s.kaiRepo.
		UpdateKaiState(
			ctx,
			tx,
			kaiState,
		)
}

func (s *Service) checkEvolution(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	DominantAttributeID uuid.UUID,
) error {
	return nil
}

func (s *Service) generateMotivationalMessage(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	dominantAttributeID *uuid.UUID,
) error {

	if dominantAttributeID == nil {
		return nil
	}

	message, err := s.messageRepo.
		FindRandomMessageByAttribute(
			ctx,
			*dominantAttributeID,
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al buscar un mensaje motivacional",
		)
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
