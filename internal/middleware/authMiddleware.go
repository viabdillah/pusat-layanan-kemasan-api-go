// internal/middleware/authMiddleware.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tidak terotorisasi, token tidak ditemukan"})
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tidak terotorisasi, token tidak valid"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userIDHex, ok := claims["id"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Klaim token tidak valid"})
				return
			}

			userID, err := primitive.ObjectIDFromHex(userIDHex)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna di token tidak valid"})
				return
			}

			var user models.User
			userCollection := config.DB.Collection("users")
			err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
				return
			}

			// Lampirkan data user ke context Gin
			c.Set("user", user)
			c.Next() // Lanjutkan ke handler berikutnya
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Klaim token tidak valid"})
		}
	}
}

func IsAdmin() gin.HandlerFunc {
	return HasRole(models.RoleAdmin)
}

// BUAT FUNGSI BARU YANG FLEKSIBEL INI
func HasRole(requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Pengguna tidak ditemukan di context"})
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Tipe data pengguna tidak valid"})
			return
		}

		// Cek apakah peran pengguna sesuai dengan yang dibutuhkan
		if user.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Peran tidak memadai."})
			return
		}

		c.Next()
	}
}

// Middleware baru untuk memeriksa salah satu dari beberapa peran
func HasAnyRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Pengguna tidak ditemukan di context"})
			return
		}

		user := userInterface.(models.User)

		isAllowed := false
		for _, role := range roles {
			if user.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Peran tidak sesuai."})
			return
		}

		c.Next()
	}
}
