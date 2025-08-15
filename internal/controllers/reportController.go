// internal/controllers/reportController.go
package controllers

import (
	"context"
	"fmt" // Tambahkan fmt jika belum ada
	"net/http"
	"time"

	"pusat-layanan-kemasan/backend-go/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2" // Impor excelize
	"go.mongodb.org/mongo-driver/bson"
)

func GetSalesSummary(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")

	// Ambil query parameter 'period', default-nya 'monthly'
	period := c.DefaultQuery("period", "monthly")

	now := time.Now()
	var startDate, endDate time.Time

	// Hitung rentang tanggal berdasarkan periode yang diminta
	switch period {
	case "daily":
		// Dari awal hari ini sampai akhir hari ini
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startDate = startOfDay
		endDate = startOfDay.Add(24*time.Hour - 1*time.Nanosecond)
	case "weekly":
		// Dari awal minggu ini (Senin) sampai akhir minggu ini (Minggu)
		weekday := now.Weekday()
		// Go menganggap Minggu = 0, Senin = 1, ... Sabtu = 6
		// Kita normalisasi agar Senin = 0
		offset := int(time.Monday - weekday)
		if weekday == time.Sunday {
			offset = -6 // Jika hari ini Minggu, anggap bagian dari minggu lalu
		}
		startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)
		startDate = startOfWeek
		endDate = startOfWeek.AddDate(0, 0, 7).Add(-1 * time.Nanosecond)
	case "yearly":
		// Dari 1 Januari tahun ini sampai akhir tahun ini
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		startDate = startOfYear
		endDate = startOfYear.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
	case "monthly":
		// Dari awal bulan ini sampai akhir bulan ini
		fallthrough // Jalankan case default
	default:
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = startOfMonth
		endDate = startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status": "selesai",
				"updatedAt": bson.M{
					"$gte": startDate, // Lebih besar atau sama dengan tanggal mulai
					"$lte": endDate,   // Lebih kecil atau sama dengan tanggal akhir
				},
			},
		},
		{
			"$group": bson.M{
				"_id":          nil,
				"totalRevenue": bson.M{"$sum": "$totalPrice"},
				"totalOrders":  bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := orderCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menjalankan agregasi"})
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil laporan"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"totalRevenue": 0,
			"totalOrders":  0,
			"period":       period,
			"startDate":    startDate,
			"endDate":      endDate,
		})
		return
	}

	// Tambahkan informasi periode ke respons
	response := results[0]
	response["period"] = period
	response["startDate"] = startDate
	response["endDate"] = endDate

	c.JSON(http.StatusOK, response)
}

// GANTI FUNGSI LAMA DENGAN VERSI BARU INI
func ExportSalesSummaryToExcel(c *gin.Context) {
	// 1. Logika Perhitungan Tanggal (Tidak berubah)
	period := c.DefaultQuery("period", "monthly")
	now := time.Now()
	var startDate, endDate time.Time
	// ... (copy-paste seluruh blok 'switch period' dari atas ke sini)
	switch period {
	case "daily":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startDate = startOfDay
		endDate = startOfDay.Add(24*time.Hour - 1*time.Nanosecond)
	case "weekly":
		weekday := now.Weekday()
		offset := int(time.Monday - weekday)
		if weekday == time.Sunday {
			offset = -6
		}
		startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)
		startDate = startOfWeek
		endDate = startOfWeek.AddDate(0, 0, 7).Add(-1 * time.Nanosecond)
	case "yearly":
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		startDate = startOfYear
		endDate = startOfYear.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
	case "monthly":
		fallthrough
	default:
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = startOfMonth
		endDate = startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
	}

	// 2. Buat Pipeline Agregasi dengan $lookup
	orderCollection := config.DB.Collection("orders")
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status":    "selesai",
				"updatedAt": bson.M{"$gte": startDate, "$lte": endDate},
			},
		},
		{
			// Tahap JOIN dengan koleksi 'customers'
			"$lookup": bson.M{
				"from":         "customers",    // Koleksi tujuan
				"localField":   "customerId",   // Field di koleksi 'orders'
				"foreignField": "_id",          // Field di koleksi 'customers'
				"as":           "customerInfo", // Nama field baru untuk hasil join
			},
		},
		{
			// $lookup menghasilkan array, kita ubah menjadi objek tunggal
			"$unwind": "$customerInfo",
		},
		{
			// Urutkan berdasarkan tanggal
			"$sort": bson.M{"updatedAt": -1},
		},
	}

	// Struct untuk menampung hasil gabungan
	var results []struct {
		OrderID      string    `bson:"orderId"`
		TotalPrice   float64   `bson:"totalPrice"`
		UpdatedAt    time.Time `bson:"updatedAt"`
		CustomerInfo struct {
			Name string `bson:"name"`
		} `bson:"customerInfo"`
	}

	cursor, err := orderCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data laporan"})
		return
	}

	// 3. Buat File Excel (Logika ini sekarang menggunakan 'results')
	f := excelize.NewFile()
	sheetName := "Laporan Penjualan"
	index, _ := f.NewSheet(sheetName)

	headers := []string{"ID Pesanan", "Nama Pelanggan", "Total Harga", "Tanggal Selesai"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, result := range results {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), result.OrderID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), result.CustomerInfo.Name) // <-- PERBAIKAN DI SINI
		f.SetCellFloat(sheetName, fmt.Sprintf("C%d", row), result.TotalPrice, 2, 64)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), result.UpdatedAt.Format("2 January 2006"))
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// 4. Kirim File sebagai Respons HTTP (Tidak berubah)
	fileName := fmt.Sprintf("Laporan-Penjualan-%s-%s.xlsx", period, now.Format("2006-01-02"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis file excel"})
	}
}

// GetSalesChartData menyiapkan data untuk ditampilkan dalam bentuk grafik
func GetSalesChartData(c *gin.Context) {
	orderCollection := config.DB.Collection("orders")
	period := c.DefaultQuery("period", "monthly")
	now := time.Now()
	var startDate, endDate time.Time

	// ... (copy-paste blok 'switch period' yang sama persis dari GetSalesSummary)
	switch period {
	case "daily":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startDate = startOfDay
		endDate = startOfDay.Add(24*time.Hour - 1*time.Nanosecond)
	case "weekly":
		weekday := now.Weekday()
		offset := int(time.Monday - weekday)
		if weekday == time.Sunday {
			offset = -6
		}
		startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)
		startDate = startOfWeek
		endDate = startOfWeek.AddDate(0, 0, 7).Add(-1 * time.Nanosecond)
	case "yearly":
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		startDate = startOfYear
		endDate = startOfYear.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
	case "monthly":
		fallthrough
	default:
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = startOfMonth
		endDate = startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status":    "selesai",
				"updatedAt": bson.M{"$gte": startDate, "$lte": endDate},
			},
		},
		{
			// Kelompokkan data berdasarkan HARI
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$updatedAt"},
				},
				"totalRevenue": bson.M{"$sum": "$totalPrice"},
			},
		},
		{
			// Urutkan berdasarkan tanggal
			"$sort": bson.M{"_id": 1},
		},
	}

	cursor, err := orderCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data grafik"})
		return
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data grafik"})
		return
	}

	c.JSON(http.StatusOK, results)
}
