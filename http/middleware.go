package http

import (
	"errors"
	jwt "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	CtxUserIDKey        = "user_id"
	CtxUserRoleKey      = "user_role"
)

type TokenValidator interface {
	ValidateAccessToken(tokenStr string) (*jwt.TokenClaims, error)
}

type AuthMiddleware struct {
	validator TokenValidator
}

func NewAuthMiddleware(validator TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, BearerPrefix))
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token string cannot be empty",
			})
			return
		}

		claims, err := m.validator.ValidateAccessToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrExpiredToken) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  "TOKEN_EXPIRED",
					"error": "access token has expired, please refresh token",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  "INVALID_TOKEN",
				"error": "unauthorized access",
			})
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUserRoleKey, claims.Role)

		c.Next()
	}
}