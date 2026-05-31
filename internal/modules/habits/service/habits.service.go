package habitsService

import (
	"context"
	"log"
	"net/http"
	"time"

	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	repositoryPort "kai-back/internal/modules/habits/repository"
	errorHandler "kai-back/internal/shared/errors"

	"github.com/google/uuid"
)

type ServicePort interface {
	GetUserHabits(ctx context.Context, userID uuid.UUID) (*habitsdto.HabitsViewResponse, error)
	GetCategories(ctx context.Context) ([]habitsdto.HabitCategoryResponse, error)
	GetCatalogByCategory(
		ctx context.Context,
		userID uuid.UUID,
		categoryID uuid.UUID,
	) ([]habitsdto.HabitCatalogResponse, error)
	SelectHabit(
		ctx context.Context,
		userID uuid.UUID,
		req habitsdto.SelectHabitRequest,
	) (*habitsdto.SelectHabitResponse, error)
}

type Service struct {
	repository repositoryPort.HabitsRepository
}

func NewService(repository repositoryPort.HabitsRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetUserHabits(ctx context.Context, userID uuid.UUID) (*habitsdto.HabitsViewResponse, error) {
	userHabits, err := s.repository.FindUserHabits(ctx, userID)
	if err != nil {
		log.Printf("Error fetching user habits: %v", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "error al obtener hábitos del usuario")
	}

	completedToday, err := s.repository.CountDailyCompleted(ctx, userID)
	if err != nil {
		log.Printf("Error counting daily completed habits: %v", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "error al contar hábitos completados hoy")
	}

	habitsResponse := make([]habitsdto.UserHabitResponse, 0, len(userHabits))
	for _, habit := range userHabits {
		habitsResponse = append(habitsResponse, buildUserHabitResponse(habit))
	}

	total := len(habitsResponse)
	if completedToday > total {
		completedToday = total
	}

	return &habitsdto.HabitsViewResponse{
		DailyProgress: habitsdto.DailyProgress{
			Total:     total,
			Completed: completedToday,
			Pending:   total - completedToday,
		},
		Habits: habitsResponse,
	}, nil
}

func buildUserHabitResponse(habit habitsmodel.UserHabit) habitsdto.UserHabitResponse {
	// records := make([]habitsdto.HabitRecordResponse, 0, len(habit.HabitRecords))
	completedToday := false
	totalXP := 0

	for _, record := range habit.HabitRecords {
		if isToday(record.Fecha) && record.Completado {
			completedToday = true
		}
		if record.Completado {
			totalXP += record.XPGanada
		}

		// records = append(records, habitsdto.HabitRecordResponse{
		// 	ID:              record.ID.String(),
		// 	Date:            record.Fecha,
		// 	Completed:       record.Completado,
		// 	RegisteredValue: record.ValorRegistrado,
		// 	XPEarned:        record.XPGanada,
		// 	CreatedAt:       record.CreatedAt,
		// })
	}

	return habitsdto.UserHabitResponse{
		ID:             habit.ID.String(),
		HabitCatalogID: habit.HabitCatalogID.String(),
		Name:           habit.HabitCatalog.Name,
		Description:    habit.HabitCatalog.Description,
		Category:       habit.HabitCatalog.Category,
		CareType:       habit.HabitCatalog.CareType,
		Difficulty:     habit.HabitCatalog.Difficulty,
		BaseXP:         habit.HabitCatalog.BaseXP,
		Custom:         habit.Personalized,
		Active:         habit.Active,
		StartDate:      habit.StartDate,
		HabitImage:     habit.HabitCatalog.HabitImage,
		CompletedToday: completedToday,
		TotalXP:        totalXP,
		CurrentStreak:  calculateCurrentStreak(habit.HabitRecords),
		// Records:        records,
	}
}

func isToday(value time.Time) bool {
	now := time.Now()
	year, month, day := now.Date()
	valueYear, valueMonth, valueDay := value.Date()

	return year == valueYear && month == valueMonth && day == valueDay
}

func calculateCurrentStreak(records []habitsmodel.HabitRecord) int {
	completedDates := make(map[string]bool)

	for _, record := range records {
		if !record.Completado {
			continue
		}

		completedDates[dateKey(record.Fecha)] = true
	}

	if len(completedDates) == 0 {
		return 0
	}

	today := dateOnly(time.Now())
	yesterday := today.AddDate(0, 0, -1)

	var currentDate time.Time
	switch {
	case completedDates[dateKey(today)]:
		currentDate = today
	case completedDates[dateKey(yesterday)]:
		currentDate = yesterday
	default:
		return 0
	}

	streak := 0
	for completedDates[dateKey(currentDate)] {
		streak++
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	return streak
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func dateKey(value time.Time) string {
	return dateOnly(value).Format("2006-01-02")
}

func (s *Service) GetCategories(
	ctx context.Context,
) ([]habitsdto.HabitCategoryResponse, error) {

	return s.repository.FindCategories(ctx)
}

func (s *Service) GetCatalogByCategory(
	ctx context.Context,
	userID uuid.UUID,
	categoryID uuid.UUID,
) ([]habitsdto.HabitCatalogResponse, error) {

	habits, err := s.repository.FindCatalogByCategoryExcludingUser(
		ctx,
		userID,
		categoryID,
	)
	if err != nil {
		return nil, err
	}

	response := make([]habitsdto.HabitCatalogResponse, 0, len(habits))

	for _, habit := range habits {

		response = append(response, habitsdto.HabitCatalogResponse{
			ID:           habit.ID,
			Name:         habit.Name,
			Description:  habit.Description,
			Category:     habit.Category,
			CareType:     habit.CareType,
			Difficulty:   habit.Difficulty,
			BaseXP:       habit.BaseXP,
			Premium:      habit.IsPremium,
			HabitImage:   habit.HabitImage,
			CategoryXPID: habit.XPCategoryID,
		})
	}

	return response, nil
}

func (s *Service) SelectHabit(
	ctx context.Context,
	userID uuid.UUID,
	req habitsdto.SelectHabitRequest,
) (*habitsdto.SelectHabitResponse, error) {

	habitCatalog, err := s.repository.FindCatalogByID(
		ctx,
		req.HabitCatalogID,
	)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"habito no encontrado",
		)
	}

	exists, err := s.repository.UserHasActiveHabit(
		ctx,
		userID,
		req.HabitCatalogID,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errorHandler.NewAppError(
			http.StatusConflict,
			"el usuario ya posee este habito",
		)
	}

	userHabit := &habitsmodel.UserHabit{
		UserID:         userID,
		HabitCatalogID: habitCatalog.ID,
		Personalized:   false,
		Active:         true,
	}

	initialRecord := &habitsmodel.HabitRecord{
		UsuarioID:  userID,
		Fecha:      time.Now(),
		Completado: false,
		XPGanada:   0,
	}

	if err := s.repository.SelectHabit(
		ctx,
		userHabit,
		initialRecord,
	); err != nil {
		return nil, err
	}

	return &habitsdto.SelectHabitResponse{
		ID:             userHabit.ID,
		HabitCatalogID: userHabit.HabitCatalogID,
		Active:         userHabit.Active,
	}, nil
}
