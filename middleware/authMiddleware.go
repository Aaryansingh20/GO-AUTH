package middleware

import (
	"net/http"
	"strings"

	helpers "github.com/Aaryansingh20/jwt/helpers"
	"github.com/gin-gonic/gin"
)
func Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
            c.Abort()
            return
        }

        // Expect header like: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        clientToken := parts[1]

        claims, err := helpers.ValidateToken(clientToken)
        if err != "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err})
            c.Abort()
            return
        }

        c.Set("email", claims.Email)
        c.Set("first_name", claims.First_name)
        c.Set("last_name", claims.Last_name)
        c.Set("uid", claims.Uid)
        c.Set("user_type", claims.User_type)

        c.Next()
    }
}
