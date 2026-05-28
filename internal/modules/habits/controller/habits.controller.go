package habitsController

import (
	"net/http"

	habitsdto "kai-back/internal/modules/habits/dto"
	servicePort "kai-back/internal/modules/habits/service"
	errorHandler "kai-back/internal/shared/errors"
	"kai-back/internal/shared/helpers"
	"kai-back/internal/shared/response"
	validatorx "kai-back/internal/shared/validator"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service servicePort.ServicePort
}

func NewController(service servicePort.ServicePort) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) GetHabits(c *gin.Context) {
	userID, err := helpers.ValidateUserUUID(c)
	if err != nil {
		return
	}

	habits, err := ctrl.service.GetUserHabits(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(habits))
}

func (ctrl *Controller) GetCategories(c *gin.Context) {

	resp, err := ctrl.service.GetCategories(
		c.Request.Context(),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(resp),
	)
}

func (ctrl *Controller) GetCatalogByCategory(c *gin.Context) {

	categoryID, err := helpers.ValidateCategoryID(c)
	if err != nil {
		return
	}

	userID, err := helpers.ValidateUserUUID(c)
	if err != nil {
		return
	}

	resp, err := ctrl.service.GetCatalogByCategory(
		c.Request.Context(),
		userID,
		categoryID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(resp),
	)
}

func (ctrl *Controller) SelectHabit(c *gin.Context) {

	var req habitsdto.SelectHabitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"request json invalido",
			),
		)
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				message,
			),
		)
		return
	}

	userID, err := helpers.ValidateUserUUID(c)
	if err != nil {
		return
	}

	resp, err := ctrl.service.SelectHabit(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusCreated,
		response.Success(resp),
	)
}
