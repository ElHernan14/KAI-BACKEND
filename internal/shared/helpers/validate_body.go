package helpers

import (
	"net/http"

	errorHandler "kai-back/internal/shared/errors"
	validatorx "kai-back/internal/shared/validator"

	"github.com/gin-gonic/gin"
)

func BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"request json invalido",
			),
		)
		return true
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				message,
			),
		)
		return true
	}

	return false
}
