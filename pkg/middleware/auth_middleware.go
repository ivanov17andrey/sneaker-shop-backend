package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"github.com/rs/zerolog/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserClaims struct {
	UserID uint
	Name   string
	jwt.RegisteredClaims
}

func ValidateToken(token string) (bool, uint) {
	tokenInfo, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Error parsing token")
		return false, 0
	}

	if claims, ok := tokenInfo.Claims.(*UserClaims); ok && tokenInfo.Valid {
		return true, claims.UserID
	}

	log.Error().Msg("Invalid token")
	return false, 0
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := authHeader[len("Bearer "):]
		valid, userID := ValidateToken(token)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
