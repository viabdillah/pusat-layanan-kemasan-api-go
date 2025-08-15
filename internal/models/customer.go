// internal/models/customer.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer mendefinisikan struktur data untuk seorang pelanggan
type Customer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	PhoneNumber string             `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`

	// Nomor-nomor registrasi produk milik pelanggan
	PirtNumber  string `bson:"pirtNumber,omitempty" json:"pirtNumber,omitempty"`
	HalalNumber string `bson:"halalNumber,omitempty" json:"halalNumber,omitempty"`
	BpomNumber  string `bson:"bpomNumber,omitempty" json:"bpomNumber,omitempty"`

	// Timestamps
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
