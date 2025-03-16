package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
    UserID   int    `json:"user_id"`
    TenantID int    `json:"tenant_id"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token
func GenerateToken(userID, tenantID int, email string) (string, error) {
    // Create claims with multiple fields
    claims := &Claims{
        UserID:   userID,
        TenantID: tenantID,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

// AuthMiddleware verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // Remove Bearer prefix if it exists
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        tokenString = strings.TrimSpace(tokenString)

        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET_KEY")), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Set user information in context
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)
        c.Set("email", claims.Email)

        c.Next()
    }
} 