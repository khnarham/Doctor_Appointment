package http

import (
	"net/http"

	jwt "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"

	"github.com/gin-gonic/gin"
)

func RequireAuth(tokenProvider jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication token missing"})
			return
		}

		claims, err := tokenProvider.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}