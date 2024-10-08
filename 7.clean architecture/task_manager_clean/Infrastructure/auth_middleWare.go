package Infrastructure

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	// "golang.org/x/crypto/bcrypt"
	"main/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation logic

		// JWT validation logic
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		// Extract the user data from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to parse token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		fmt.Println(userID)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to parse user ID"})
			c.Abort()
			return
		}
		IsAdmin, ok := claims["is_admin"].(bool)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to parse user ID"})
			c.Abort()
			return

		}
		
		c.Set("user_id", userID)
		c.Set("is_admin", IsAdmin)

		c.Next()
	}
}

func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// userID,err := c.Set("user_id" )
		IsAdmin, ok := c.Get("is_admin")

		if ok == true && IsAdmin.(bool) == true {
			c.Next()
		}

	}
}
