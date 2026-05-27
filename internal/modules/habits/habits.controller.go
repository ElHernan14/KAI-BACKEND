package habits

import (
	"net/http"

	appcontext "kai-back/internal/shared/context"
	errorHandler "kai-back/internal/shared/errors"
	"kai-back/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	service ServicePort
}

func NewController(service ServicePort) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) GetHabits(c *gin.Context) {
	userIDString, exists := appcontext.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(errorHandler.NewAppError(http.StatusUnauthorized, "usuario no autenticado"))
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(errorHandler.NewAppError(http.StatusUnauthorized, "usuario no autenticado"))
		return
	}

	habits, err := ctrl.service.GetUserHabits(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(habits))
}
