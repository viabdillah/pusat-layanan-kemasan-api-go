// internal/controllers/customerController.go
package controllers

import (
	"context"
	"net/http"
	"time"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCustomer membuat pelanggan baru
func CreateCustomer(c *gin.Context) {
	var input models.Customer
	customerCollection := config.DB.Collection("customers")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	// Cek duplikasi berdasarkan nama pelanggan
	var existingCustomer models.Customer
	err := customerCollection.FindOne(context.Background(), bson.M{"name": input.Name}).Decode(&existingCustomer)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pelanggan dengan nama ini sudah ada"})
		return
	}

	newCustomer := models.Customer{
		ID:          primitive.NewObjectID(),
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		PirtNumber:  input.PirtNumber,
		HalalNumber: input.HalalNumber,
		BpomNumber:  input.BpomNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = customerCollection.InsertOne(context.Background(), newCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pelanggan baru"})
		return
	}

	c.JSON(http.StatusCreated, newCustomer)
}

// GetCustomers mengambil daftar pelanggan, dengan opsi pencarian
func GetCustomers(c *gin.Context) {
	customerCollection := config.DB.Collection("customers")

	// Ambil query parameter 'search' dari URL, contoh: /api/customers?search=Toko
	searchQuery := c.Query("search")
	filter := bson.M{}

	if searchQuery != "" {
		// Buat filter pencarian case-insensitive berdasarkan nama
		filter["name"] = bson.M{"$regex": searchQuery, "$options": "i"}
	}

	var customers []models.Customer
	cursor, err := customerCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan"})
		return
	}

	if err = cursor.All(context.Background(), &customers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pelanggan"})
		return
	}

	// Jika tidak ada pelanggan ditemukan, kembalikan array kosong, bukan null
	if customers == nil {
		customers = []models.Customer{}
	}

	c.JSON(http.StatusOK, customers)
}
