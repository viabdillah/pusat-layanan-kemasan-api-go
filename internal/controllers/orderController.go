// internal/controllers/orderController.go
package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(c *gin.Context) {
	// Dapatkan pengguna (kasir) yang sedang login
	userInterface, _ := c.Get("user")
	kasir := userInterface.(models.User)
	orderCollection := config.DB.Collection("orders")

	// Definisikan struct untuk menampung input dari JSON
	var input struct {
		CustomerID string             `json:"customerId" binding:"required"`
		Items      []models.OrderItem `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	// Konversi CustomerID dari string ke ObjectID
	customerID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CustomerID tidak valid"})
		return
	}

	// Hitung TotalPrice dan tentukan Status Awal dari server
	var totalPrice float64
	// Status awal diasumsikan menunggu produksi, kecuali ada item yang butuh desain
	initialStatus := models.StatusMenungguProduksi

	for _, item := range input.Items {
		totalPrice += item.PricePerPiece * float64(item.Quantity)
		if !item.HasDesign {
			initialStatus = models.StatusMenungguDesain // Jika ada 1 saja item butuh desain, seluruh order menunggu desain
		}
	}

	// Buat objek pesanan baru dengan struktur yang benar
	newOrder := models.Order{
		ID:         primitive.NewObjectID(),
		OrderID:    fmt.Sprintf("ORD-%s", time.Now().Format("20060102-150405")),
		CustomerID: customerID,
		Items:      input.Items,
		TotalPrice: totalPrice,
		Status:     initialStatus,
		KasirID:    kasir.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Simpan pesanan baru ke database
	_, err = orderCollection.InsertOne(context.Background(), newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesanan"})
		return
	}

	c.JSON(http.StatusCreated, newOrder)
}

// GetDesignerQueue mengambil daftar pekerjaan untuk desainer
func GetDesignerQueue(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")
	user := c.MustGet("user").(models.User) // Dapatkan user dari context

	// Filter: tampilkan pesanan yang statusnya 'menunggu_desain' ATAU
	// yang sudah ditugaskan ke desainer yang sedang login ini.
	filter := bson.M{
		"$or": []bson.M{
			{"status": models.StatusMenungguDesain},
			{"designerId": user.ID, "status": models.StatusProsesDesain},
		},
	}

	var orders []models.Order
	cursor, err := orderCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil antrian desainer"})
		return
	}

	if err = cursor.All(context.Background(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses antrian desainer"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOperatorQueue mengambil daftar pekerjaan untuk operator
func GetOperatorQueue(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")
	user := c.MustGet("user").(models.User)

	// Filter: tampilkan pesanan yang statusnya 'menunggu_produksi' ATAU
	// yang sudah ditugaskan ke operator yang sedang login ini.
	filter := bson.M{
		"$or": []bson.M{
			{"status": models.StatusMenungguProduksi},
			{"operatorId": user.ID, "status": models.StatusProsesProduksi},
		},
	}

	var orders []models.Order
	cursor, err := orderCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil antrian operator"})
		return
	}

	if err = cursor.All(context.Background(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses antrian operator"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetKasirQueue mengambil daftar pekerjaan untuk kasir (pesanan yang siap bayar)
func GetKasirQueue(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")

	// Filter: tampilkan pesanan yang statusnya 'menunggu_pembayaran'
	filter := bson.M{"status": models.StatusMenungguPembayaran}

	var orders []models.Order
	cursor, err := orderCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil antrian kasir"})
		return
	}

	if err = cursor.All(context.Background(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses antrian kasir"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus memperbarui status sebuah pesanan dengan validasi peran
func UpdateOrderStatus(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")
	user := c.MustGet("user").(models.User)

	orderID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pesanan tidak valid"})
		return
	}

	var input struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status baru dibutuhkan"})
		return
	}

	// Tentukan update yang akan dilakukan
	updateSet := bson.M{
		"status":    input.Status,
		"updatedAt": time.Now(),
	}

	// Logika spesifik berdasarkan peran
	switch user.Role {
	case models.RoleDesigner:
		// Desainer mulai mengerjakan pesanan
		if input.Status == models.StatusProsesDesain {
			updateSet["designerId"] = user.ID
		}
	case models.RoleOperator:
		// Operator mulai mengerjakan pesanan
		if input.Status == models.StatusProsesProduksi {
			updateSet["operatorId"] = user.ID
		}
	case models.RoleKasir:
		// Kasir hanya boleh mengubah status menjadi 'selesai'
		if input.Status == models.StatusSelesai {
			// Tidak ada field tambahan yang di-set, hanya status
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Kasir hanya dapat menyelesaikan pesanan"})
			return
		}

	case models.RoleAdmin:
		// Admin bisa melakukan apa saja, tidak ada tambahan logika
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengubah status pesanan"})
		return
	}

	update := bson.M{"$set": updateSet}

	result, err := orderCollection.UpdateOne(context.Background(), bson.M{"_id": orderID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status pesanan"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pesanan berhasil diperbarui"})
}

// Di file: internal/controllers/orderController.go

// GANTI FUNGSI LAMA DENGAN VERSI BARU INI
func GetMonitoringOrders(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")

	// Pipeline untuk mengambil pesanan aktif DAN menggabungkannya dengan data pelanggan
	pipeline := []bson.M{
		{
			// Tahap 1: Filter pesanan yang aktif
			"$match": bson.M{
				"status": bson.M{
					"$nin": []models.OrderStatus{models.StatusSelesai, models.StatusDibatalkan},
				},
			},
		},
		{
			// Tahap 2: JOIN dengan koleksi 'customers'
			"$lookup": bson.M{
				"from":         "customers",
				"localField":   "customerId",
				"foreignField": "_id",
				"as":           "customerInfo",
			},
		},
		{
			// Ubah 'customerInfo' dari array menjadi objek
			"$unwind": bson.M{
				"path":                       "$customerInfo",
				"preserveNullAndEmptyArrays": true, // Jaga pesanan tetap ada meskipun pelanggan terhapus
			},
		},
		{
			// Urutkan berdasarkan tanggal update terakhir
			"$sort": bson.M{"updatedAt": -1},
		},
	}

	// Struct untuk menampung hasil gabungan
	var results []bson.M // Kita gunakan bson.M agar lebih fleksibel

	cursor, err := orderCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pemantauan"})
		return
	}

	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pemantauan"})
		return
	}

	if results == nil {
		results = []bson.M{}
	}

	c.JSON(http.StatusOK, results)
}

// GetOrderByID mengambil satu pesanan berdasarkan ID-nya
func GetOrderByID(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")
	orderID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pesanan tidak valid"})
		return
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": orderID},
		},
		{
			"$lookup": bson.M{
				"from":         "customers",
				"localField":   "customerId",
				"foreignField": "_id",
				"as":           "customerInfo",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$customerInfo",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	var results []bson.M
	cursor, err := orderCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}

	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pesanan"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, results[0])
}
