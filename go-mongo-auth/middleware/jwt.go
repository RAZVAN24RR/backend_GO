package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// jwtSecret trebuie să fie aceeași cu cea din controllers.
var jwtSecret = []byte("secret-key")

// AuthMiddleware validează token-ul JWT preluat din header-ul Authorization.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header-ul Authorization lipsește"})
			return
		}

		// Se așteaptă formatul: "Bearer <token>"
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formatul token-ului este invalid"})
			return
		}

		tokenString := fields[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificăm metoda de semnare
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token-ul nu conține claim-uri valide"})
			return
		}

		// Stocăm claim-urile în context pentru handler-ele ulterioare.
		c.Set("user", claims)
		c.Next()
	}
}