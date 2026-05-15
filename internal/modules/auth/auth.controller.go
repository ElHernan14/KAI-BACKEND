package auth

import (
	"context"
	"net/http"

	authdto "kai-back/internal/modules/auth/dto"
	errorHandler "kai-back/internal/shared/errors"
	"kai-back/internal/shared/response"
	validatorx "kai-back/internal/shared/validator"

	"github.com/gin-gonic/gin"
)

type ServicePort interface {
	Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.AuthResponse, error)
	Login(ctx context.Context, req authdto.LoginRequest) (*authdto.AuthResponse, error)
}

type Controller struct {
	service ServicePort
}

func NewController(service ServicePort) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errorHandler.NewAppError(http.StatusBadRequest, "request json invalido"))
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(errorHandler.NewAppError(http.StatusBadRequest, message))
		return
	}

	authResponse, err := ctrl.service.Register(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessWithCode(http.StatusCreated, authResponse))
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errorHandler.NewAppError(http.StatusBadRequest, "request json invalido"))
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(errorHandler.NewAppError(http.StatusBadRequest, message))
		return
	}

	authResponse, err := ctrl.service.Login(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(authResponse))
}
