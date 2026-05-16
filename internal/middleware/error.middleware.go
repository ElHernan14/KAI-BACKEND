package middleware

import (
	"errors"
	"log"
	"net/http"

	errorHandler "kai-back/internal/shared/errors"
	"kai-back/internal/shared/response"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID := GetRequestID(c)
				log.Printf("ERROR middleware | request_id=%s | panic=%v", requestID, recovered)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.Error(http.StatusInternalServerError, errorHandler.ErrInternal.Error()),
				)
			}
		}()

		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		requestID := GetRequestID(c)

		var appErr errorHandler.AppError
		if errors.As(err, &appErr) {
			log.Printf("ERROR | request_id=%s | status=%d | error=%s", requestID, appErr.Code, appErr.Message)
			c.AbortWithStatusJSON(appErr.Code, response.Error(appErr.Code, appErr.Message))
			return
		}

		log.Printf("ERROR middleware | request_id=%s | unexpected_error=%v", requestID, err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.Error(http.StatusInternalServerError, errorHandler.ErrInternal.Error()),
		)
	}
}
