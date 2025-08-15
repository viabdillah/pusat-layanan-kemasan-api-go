// cmd/server/main.go
package main

import (
	"log"
	"os"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/routes" // Kita akan uncomment ini nanti

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// fungsi init() akan berjalan sebelum main()
func init() {
	// Muat file .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan")
	}
	// Hubungkan ke DB
	config.ConnectDB()
}

func main() {
	router := gin.Default()

	// --- 2. KONFIGURASI DAN TERAPKAN MIDDLEWARE CORS ---
	// Konfigurasi ini memberitahu backend kita untuk mengizinkan
	// permintaan dari server frontend Vue kita.
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Izinkan origin frontend
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config)) // Terapkan middleware ke semua rute

	// ... rute tes bisa dihapus atau dibiarkan ...
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Pusat Layanan Kemasan dengan Golang sedang berjalan...",
		})
	})

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
