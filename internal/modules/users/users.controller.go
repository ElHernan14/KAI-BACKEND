package users

import (
	"context"
	usersdto "kai-back/internal/modules/users/dto"
	contextutil "kai-back/internal/shared/context"
	errorHandler "kai-back/internal/shared/errors"
	"kai-back/internal/shared/response"
	validatorx "kai-back/internal/shared/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServicePort interface {
	GetMe(ctx context.Context, userID uuid.UUID) (*usersdto.MeResponse, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req usersdto.UpdateMeRequest) (*usersdto.MeResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req usersdto.ChangePasswordRequest) error
}

type Controller struct {
	service ServicePort
}

func NewController(service ServicePort) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Me(c *gin.Context) {

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticado",
			),
		)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := ctrl.service.GetMe(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(user),
	)
}

func (ctrl *Controller) UpdateMe(c *gin.Context) {
	var req usersdto.UpdateMeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errorHandler.NewAppError(
			http.StatusBadRequest,
			"request json invalido",
		))
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(errorHandler.NewAppError(
			http.StatusBadRequest,
			message,
		))
		return
	}

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticados",
			),
		)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(err)
		return
	}

	responseData, err := ctrl.service.UpdateMe(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(responseData),
	)
}

func (ctrl *Controller) ChangePassword(c *gin.Context) {

	var req usersdto.ChangePasswordRequest

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

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticado",
			),
		)
		return
	}

	userID, errParse := uuid.Parse(userIDString)
	if errParse != nil {
		_ = c.Error(errParse)
		return
	}

	err := ctrl.service.ChangePassword(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.SuccessWithCode(
			http.StatusOK,
			"contraseña actualizada correctamente",
		),
	)
}
