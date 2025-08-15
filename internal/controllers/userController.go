// internal/controllers/userController.go
package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Helper function untuk generate token JWT
func generateToken(userID primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Token berlaku 30 hari
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func Register(c *gin.Context) {
	var user models.User
	userCollection := config.DB.Collection("users")

	// 1. Bind JSON dari request body ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	// 2. Cek apakah email sudah ada
	var existingUser models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// Jika err nil, berarti user ditemukan (sudah ada)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Handle error database lainnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa email"})
		return
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}
	user.Password = string(hashedPassword)

	// 4. Set nilai default dan timestamp
	user.ID = primitive.NewObjectID()
	user.Role = models.RoleKasir
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 5. Masukkan user baru ke database
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan pengguna"})
		return
	}

	// 6. Generate JWT Token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 7. Kirim respons sukses
	c.JSON(http.StatusCreated, gin.H{
		"_id":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"token": token,
	})
}

// ... (fungsi generateToken dan Register sudah ada di atas)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	userCollection := config.DB.Collection("users")
	var user models.User

	// 1. Bind JSON dari request body ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password dibutuhkan"})
		return
	}

	// 2. Cari pengguna berdasarkan email
	err := userCollection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		// Jika user tidak ditemukan, kirim error 'Unauthorized'
		// Ini adalah praktik keamanan yang baik untuk tidak memberitahu apakah email atau password yang salah
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari pengguna"})
		return
	}

	// 3. Bandingkan password yang di-hash di DB dengan password dari input
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// Jika password tidak cocok, kirim error yang sama
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 4. Jika kredensial valid, generate token baru
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 5. Kirim respons sukses
	c.JSON(http.StatusOK, gin.H{
		"_id":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"token": token,
	})
}

func GetProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan di context"})
		return
	}

	user := userInterface.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"_id":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func GetUsers(c *gin.Context) {
	userCollection := config.DB.Collection("users")
	var users []models.User

	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pengguna"})
		return
	}

	c.JSON(http.StatusOK, users)
}
