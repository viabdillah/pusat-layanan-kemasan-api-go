// internal/models/order.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// (Tipe dan konstanta OrderStatus tidak berubah)
type OrderStatus string

const (
	StatusMenungguDesain     OrderStatus = "menunggu_desain"
	StatusProsesDesain       OrderStatus = "proses_desain"
	StatusMenungguProduksi   OrderStatus = "menunggu_produksi"
	StatusProsesProduksi     OrderStatus = "proses_produksi"
	StatusMenungguPembayaran OrderStatus = "menunggu_pembayaran"
	StatusSelesai            OrderStatus = "selesai"
	StatusDibatalkan         OrderStatus = "dibatalkan"
)

// --- STRUCT BARU ---
// OrderItem merepresentasikan satu baris produk di dalam sebuah pesanan
type OrderItem struct {
	ProductName   string  `bson:"productName" json:"productName"`
	Brand         string  `bson:"brand,omitempty" json:"brand,omitempty"` // Merk
	Label         string  `bson:"label,omitempty" json:"label,omitempty"`
	Size          string  `bson:"size,omitempty" json:"size,omitempty"`
	Quantity      int     `bson:"quantity" json:"quantity"`
	PricePerPiece float64 `bson:"pricePerPiece" json:"pricePerPiece"`
	HasDesign     bool    `bson:"hasDesign" json:"hasDesign"`
}

// --- STRUCT ORDER YANG DIPERBARUI ---
// Order mendefinisikan struktur untuk setiap pesanan
type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID    string             `bson:"orderId" json:"orderId"`
	CustomerID primitive.ObjectID `bson:"customerId" json:"customerId"` // <-- BERUBAH: Sekarang menunjuk ke ID pelanggan
	Items      []OrderItem        `bson:"items" json:"items"`           // <-- BERUBAH: Sekarang berisi daftar produk
	TotalPrice float64            `bson:"totalPrice" json:"totalPrice"` // Nilai ini akan dihitung oleh server
	Status     OrderStatus        `bson:"status" json:"status"`

	// ID Pengguna yang terlibat dalam pesanan
	KasirID    primitive.ObjectID `bson:"kasirId" json:"kasirId"`
	DesignerID primitive.ObjectID `bson:"designerId,omitempty" json:"designerId,omitempty"`
	OperatorID primitive.ObjectID `bson:"operatorId,omitempty" json:"operatorId,omitempty"`

	// Timestamps
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
