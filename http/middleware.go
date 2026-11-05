package http

import (
	"errors"
	"fmt"
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

		// DEBUG: Incoming Request details print karo
		fmt.Printf("=== [DEBUG AUTH] Path: %s | Method: %s ===\n", c.Request.URL.Path, c.Request.Method)
		fmt.Println("[DEBUG AUTH] Raw Cookie Header:", c.GetHeader("Cookie"))

		// 1. Direct Gin c.Cookie try karo
		if cookieToken, err := c.Cookie(AccessTokenCookie); err == nil && cookieToken != "" {
			tokenStr = cookieToken
			fmt.Println("[DEBUG AUTH] Token found via c.Cookie()")
		}

		// 2. Fallback: Raw Cookie header parsing
		if tokenStr == "" {
			rawCookieHeader := c.GetHeader("Cookie")
			if rawCookieHeader != "" {
				cookies := strings.Split(rawCookieHeader, ";")
				for _, cookie := range cookies {
					parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
					if len(parts) == 2 && parts[0] == AccessTokenCookie {
						tokenStr = parts[1]
						fmt.Println("[DEBUG AUTH] Token found via manual raw cookie parsing")
						break
					}
				}
			}
		}

		// 3. Fallback: Authorization Header
		if tokenStr == "" {
			authHeader := c.GetHeader(AuthorizationHeader)
			if authHeader != "" && strings.HasPrefix(authHeader, BearerPrefix) {
				tokenStr = strings.TrimSpace(strings.TrimPrefix(authHeader, BearerPrefix))
				fmt.Println("[DEBUG AUTH] Token found via Authorization Bearer Header")
			}
		}

		if tokenStr == "" {
			fmt.Println("[DEBUG AUTH ERROR] No token found anywhere in request!")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  "UNAUTHORIZED",
				"error": "authentication token is missing",
			})
			return
		}

		// Validate Token
		claims, err := m.validator.ValidateAccessToken(tokenStr)
		if err != nil {
			fmt.Printf("[DEBUG AUTH ERROR] Token Validation Failed: %v\n", err)
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

		fmt.Printf("[DEBUG AUTH SUCCESS] User ID: %s | Role: %s\n", claims.UserID, claims.Role)

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUserRoleKey, string(claims.Role)) 

		c.Next()
	}
}