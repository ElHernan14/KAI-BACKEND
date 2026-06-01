package helpers

import (
	appcontext "kai-back/internal/shared/context"
	errorHandler "kai-back/internal/shared/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidateUserUUID(c *gin.Context) (uuid.UUID, error) {
	userIDString, exists := appcontext.GetUserID(c.Request.Context())
	if !exists {
		return uuid.Nil, c.Error(errorHandler.NewAppError(http.StatusUnauthorized, "usuario no autenticado"))

	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, c.Error(errorHandler.NewAppError(http.StatusUnauthorized, "usuario no autenticado"))
	}
	return userID, nil
}

func ValidateCategoryID(c *gin.Context) (uuid.UUID, error) {
	categoryIDParam := c.Query("category_id")

	if categoryIDParam == "" {
		return uuid.Nil, c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"category_id es requerido",
			),
		)
	}

	categoryID, err := uuid.Parse(categoryIDParam)
	if err != nil {
		return uuid.Nil, c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"category_id invalido",
			),
		)
	}

	return categoryID, nil
}

func ValidateUserHabitID(c *gin.Context) (uuid.UUID, error) {
	habitUserIDParam := c.Param("habitUserId")
	if habitUserIDParam == "" {
		return uuid.Nil, c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"habit_user_id es requerido",
			),
		)
	}

	habitUserID, err := uuid.Parse(habitUserIDParam)
	if err != nil {
		return uuid.Nil, c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"habit_user_id invalido",
			),
		)
	}

	return habitUserID, nil
}
