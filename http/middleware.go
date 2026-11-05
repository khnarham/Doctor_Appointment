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
	BearerPrefix        = "Bearer"
	AccessTokenCookie   = "access_token"
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
		var tokenStr string

		// 1. Pehle direct Gin ka c.Cookie try karo
		if cookieToken, err := c.Cookie(AccessTokenCookie); err == nil && cookieToken != "" {
			tokenStr = cookieToken
		}

		// 2. Fallback: Agar c.Cookie() fail ho jaye, toh raw Cookie header se manually parse karo (Cross-domain safety)
		if tokenStr == "" {
			rawCookieHeader := c.GetHeader("Cookie")
			if rawCookieHeader != "" {
				cookies := strings.Split(rawCookieHeader, ";")
				for _, cookie := range cookies {
					parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
					if len(parts) == 2 && parts[0] == AccessTokenCookie {
						tokenStr = parts[1]
						break
					}
				}
			}
		}

		// 3. Fallback: Authorization Bearer Header check karo
		if tokenStr == "" {
			authHeader := c.GetHeader(AuthorizationHeader)
			if authHeader != "" && strings.HasPrefix(authHeader, BearerPrefix) {
				tokenStr = strings.TrimSpace(strings.TrimPrefix(authHeader, BearerPrefix))
			}
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  "UNAUTHORIZED",
				"error": "authentication token is missing",
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
		c.Set(CtxUserRoleKey, string(claims.Role)) 

		c.Next()
	}
}