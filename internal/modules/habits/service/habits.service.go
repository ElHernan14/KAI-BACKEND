package habitsService

import (
	"context"
	"log"
	"net/http"
	"time"

	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	repositoryPort "kai-back/internal/modules/habits/repository"
	helperPrivate "kai-back/internal/modules/habits/service/private"
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
	GetHabitDetail(
		ctx context.Context,
		userID uuid.UUID,
		habitUserID uuid.UUID,
	) (*habitsdto.HabitDetailResponse, error)
	DeactivateHabit(
		ctx context.Context,
		userID uuid.UUID,
		habitID uuid.UUID,
	) (*habitsdto.DeactivateHabitResponse, error)
}

type Service struct {
	repository                repositoryPort.HabitsRepository
	habitsDailyRecordsService HabitsDailyRecordsServicePort
}

func NewService(repository repositoryPort.HabitsRepository, habitsDailyRecordsService HabitsDailyRecordsServicePort) *Service {
	return &Service{
		repository:                repository,
		habitsDailyRecordsService: habitsDailyRecordsService,
	}
}

func (s *Service) GetUserHabits(ctx context.Context, userID uuid.UUID) (*habitsdto.HabitsViewResponse, error) {
	// Aseguramos que existan registros de hábitos para hoy antes de obtener los hábitos del usuario
	err := s.habitsDailyRecordsService.EnsureTodayHabitRecords(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

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
		habitsResponse = append(habitsResponse, helperPrivate.BuildUserHabitResponse(habit))
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

	//Comprobamos si el usuario ya tiene este hábito activo o inactivo para evitar duplicados y reactivar si es necesario
	existingHabit, err := s.repository.FindUserHabitByCatalogID(
		ctx,
		userID,
		req.HabitCatalogID,
	)

	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al verificar hábito existente",
		)
	}

	if existingHabit != nil {
		// Si el hábito existe y está activo, retornamos un error de conflicto
		if existingHabit.Active {

			return nil, errorHandler.NewAppError(
				http.StatusConflict,
				"el usuario ya posee este habito",
			)
		}

		// Si el hábito existe pero está inactivo, lo reactivamos
		if err := s.repository.
			ReactivateHabitWithTodayRecord(
				ctx,
				existingHabit,
			); err != nil {

			return nil, err
		}

		return &habitsdto.SelectHabitResponse{
			ID:             existingHabit.ID,
			HabitCatalogID: existingHabit.HabitCatalogID,
			Active:         true,
		}, nil
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

func (s *Service) GetHabitDetail(
	ctx context.Context,
	userID uuid.UUID,
	userHabitID uuid.UUID,
) (*habitsdto.HabitDetailResponse, error) {
	// Aseguramos que existan registros de hábitos para hoy antes de obtener los hábitos del usuario
	err := s.habitsDailyRecordsService.EnsureTodayHabitRecords(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	habit, err := s.repository.FindHabitDetailByID(
		ctx,
		userID,
		userHabitID,
	)
	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"habito del usuario no encontrado",
		)
	}

	totalXP := helperPrivate.CalculateTotalXP(habit.HabitRecords)

	currentStreak := helperPrivate.CalculateCurrentStreak(habit.HabitRecords)

	records := make([]habitsdto.HabitRecordResponse, 0)

	for _, record := range habit.HabitRecords {

		records = append(records, habitsdto.HabitRecordResponse{
			ID:              record.ID.String(),
			UserHabitID:     record.UserHabitID,
			Date:            record.Fecha,
			Completed:       record.Completado,
			RegisteredValue: record.ValorRegistrado,
			XPEarned:        record.XPGanada,
			CreatedAt:       record.CreatedAt,
		})
	}

	response := &habitsdto.HabitDetailResponse{
		Habito: habitsdto.UserHabitResponse{
			ID:             habit.ID.String(),
			HabitCatalogID: habit.HabitCatalogID.String(),

			Name:        habit.HabitCatalog.Name,
			Description: habit.HabitCatalog.Description,
			Category:    habit.HabitCatalog.Category,
			CareType:    habit.HabitCatalog.CareType,
			Difficulty:  habit.HabitCatalog.Difficulty,
			BaseXP:      habit.HabitCatalog.BaseXP,
			HabitImage:  habit.HabitCatalog.HabitImage,

			Custom:    habit.Personalized,
			Active:    habit.Active,
			StartDate: habit.StartDate,

			CompletedToday: false, // lo calculamos abajo

			TotalXP:       totalXP,
			CurrentStreak: currentStreak,
		},

		Registros: records,
	}

	today := time.Now().Format("2006-01-02")

	for _, record := range habit.HabitRecords {

		if record.Completado &&
			record.Fecha.Format("2006-01-02") == today {

			response.Habito.CompletedToday = true
			break
		}
	}

	return response, nil
}

func (s *Service) DeactivateHabit(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
) (*habitsdto.DeactivateHabitResponse, error) {

	userHabit, err := s.repository.FindUserHabitByID(
		ctx,
		userID,
		habitID,
	)

	if err != nil {
		return nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"habito no encontrado",
		)
	}

	if !userHabit.Active {
		return nil, errorHandler.NewAppError(
			http.StatusConflict,
			"el habito ya se encuentra desactivado",
		)
	}

	if err := s.repository.DeactivateHabit(
		ctx,
		habitID,
	); err != nil {
		return nil, err
	}

	return &habitsdto.DeactivateHabitResponse{
		HabitoUsuarioID: habitID.String(),
		Activo:          false,
		Mensaje:         "habito desactivado correctamente",
	}, nil
}
