package middleware

import (
	"log"
	"time"

	appcontext "kai-back/internal/shared/context"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		requestID := GetRequestID(c)
		status := c.Writer.Status()
		latency := time.Since(start)

		log.Printf(
			"request_id=%s method=%s path=%s status=%d latency=%.3fs client_ip=%s",
			requestID,
			c.Request.Method,
			path,
			status,
			latency.Seconds(),
			c.ClientIP(),
		)
		log.Println("----------------------------------------")
	}
}

func GetRequestID(c *gin.Context) string {
	if requestID, ok := appcontext.GetRequestID(c.Request.Context()); ok {
		return requestID
	}

	if requestID, ok := c.Get(string(appcontext.RequestIDKey)); ok {
		if value, ok := requestID.(string); ok {
			return value
		}
	}

	return ""
}
