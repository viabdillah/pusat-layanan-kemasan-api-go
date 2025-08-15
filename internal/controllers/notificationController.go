// internal/controllers/notificationController.go
package controllers

import (
	"context"
	"net/http"

	"pusat-layanan-kemasan/backend-go/internal/config"
	"pusat-layanan-kemasan/backend-go/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetNotifications mengambil notifikasi untuk pengguna yang sedang login
func GetNotifications(c *gin.Context) {
	notificationCollection := config.DB.Collection("notifications")
	user := c.MustGet("user").(models.User)

	// Filter notifikasi hanya untuk user ID yang sedang login
	filter := bson.M{"userId": user.ID}
	// Urutkan dari yang terbaru, batasi 20 notifikasi terakhir
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(20)

	var notifications []models.Notification
	cursor, err := notificationCollection.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}

	if err = cursor.All(context.Background(), &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses notifikasi"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead menandai notifikasi sebagai sudah dibaca
func MarkAsRead(c *gin.Context) {
	notificationCollection := config.DB.Collection("notifications")
	user := c.MustGet("user").(models.User)

	notifID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Notifikasi tidak valid"})
		return
	}

	// Filter berdasarkan ID notifikasi DAN ID pengguna
	// Ini untuk memastikan pengguna tidak bisa mengubah notifikasi orang lain
	filter := bson.M{"_id": notifID, "userId": user.ID}
	update := bson.M{"$set": bson.M{"isRead": true}}

	result, err := notificationCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui notifikasi"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan atau Anda tidak punya akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi ditandai sebagai sudah dibaca"})
}
