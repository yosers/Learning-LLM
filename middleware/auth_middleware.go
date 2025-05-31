package middleware

import (
	"net/http"
	"shofy/utils/jwt"
	"shofy/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		println(">>> AuthMiddleware called")

		var tokenString string
		var tokenFound bool

		// Get token from cookie
		if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
			tokenString = cookie
			tokenFound = true
		}
		println(">>> Token called", tokenString)

		// If no token in cookie, check Authorization header
		if !tokenFound {
			println(">>> Header called", c.GetHeader("Authorization"))

			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// Check if the header starts with "Bearer "
				splitToken := strings.Split(authHeader, "Bearer ")
				if len(splitToken) == 2 {
					tokenString = splitToken[1]
					tokenFound = true
				}
			}
		}

		// If no token found in either cookie or header
		if !tokenFound {
			println(">>> Token not Found", c.GetHeader("No authentication token provided"))

			response.Error(c, http.StatusUnauthorized, "No authentication token provided")
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user ID in context for later use
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
