package middleware

import (
	"net/http"
	"strings"

	authcore "kai-back/internal/shared/auth"
	appcontext "kai-back/internal/shared/context"
	"kai-back/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(c, "Header de autorización faltante")
			return
		}

		tokenString, ok := extractBearerToken(authHeader)
		if !ok {
			abortUnauthorized(c, "Header de autorización debe usar token Bearer")
			return
		}

		claims := &authcore.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			abortUnauthorized(c, "Token inválido o expirado")
			return
		}

		SetAuthUser(c, appcontext.AuthUser{
			UserID: claims.UserID,
			Email:  claims.Email,
		})

		c.Next()
	}
}

func extractBearerToken(authHeader string) (string, bool) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	return parts[1], true
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		response.Error(http.StatusUnauthorized, message),
	)
}

func SetAuthUser(c *gin.Context, user appcontext.AuthUser) {
	c.Set(string(appcontext.AuthUserKey), user)

	ctx := appcontext.WithAuthUser(c.Request.Context(), user)
	c.Request = c.Request.WithContext(ctx)
}

func GetAuthUser(c *gin.Context) (appcontext.AuthUser, bool) {
	if user, ok := appcontext.GetAuthUser(c.Request.Context()); ok {
		return user, true
	}

	value, ok := c.Get(string(appcontext.AuthUserKey))
	if !ok {
		return appcontext.AuthUser{}, false
	}

	user, ok := value.(appcontext.AuthUser)
	return user, ok
}
